package Lendme

//reference: https://github.com/segmentio/kafka-go
//reference: https://pkg.go.dev/github.com/segmentio/kafka-go#section-readme
import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type KafkaClient struct {
	Writer *kafka.Writer
}

// ***********************************************************************************************
// Kafka initialization
// ***********************************************************************************************
// Build a client and fail fast if config is invalid.
// Returning (*KafkaClient, error) is production-grade; callers must handle failure.
func NewKafkaClient() (*KafkaClient, error) {
	brokersCSV := strings.TrimSpace(Configuration.Kafka_Events.KafkaBrokerUrls)
	clientID := strings.TrimSpace(Configuration.Kafka_Events.KafkaClientId)

	if brokersCSV == "" {
		return nil, errors.New("kafka brokers are empty")
	}
	if clientID == "" {
		return nil, errors.New("kafka clientId is empty")
	}

	brokers := splitAndTrimCSV(brokersCSV)
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers list is empty after parsing")
	}

	c := &KafkaClient{}
	if err := c.createKafkaWriter(brokers, clientID); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *KafkaClient) Close() error {
	if c == nil || c.Writer == nil {
		return nil
	}
	return c.Writer.Close()
}

func (c *KafkaClient) createKafkaWriter(brokers []string, clientID string) error {
	log.Println("creating kafka writer")

	// Dialer: make it ready for TLS/SASL (telecom environments usually require it).
	// TODO: wire TLS/SASL config from your Configuration.
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		ClientID:  clientID,
		// TLS:         tlsConfig,
		// SASLMechanism: saslMechanism,
	}

	// IMPORTANT: if you also implement retry in Push(), keep Writer.MaxAttempts low
	// to avoid double retry storms. Otherwise, let Writer handle retries and keep Push simple.
	cfg := kafka.WriterConfig{
		Brokers: brokers,
		//Balancer: &kafka.LeastBytes{}, //if no key specified when pushing the record to kafka
		Balancer: &kafka.Hash{}, // key-based partitioning

		Dialer: dialer,

		// Timeouts: should be aligned with your Push() attempt timeouts.
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		CompressionCodec: snappy.NewCompressionCodec(),

		// Batching: good throughput defaults; tune via config for your traffic profile.
		BatchSize:    250,
		BatchBytes:   10 * 1024 * 1024, // 10MB
		BatchTimeout: 50 * time.Millisecond,

		// Delivery durability:
		// - RequireAll is safer (acks from ISR), higher latency.
		// - RequireOne is faster, less durable.
		RequiredAcks: -1, // RequireAll (production-safe)

		// Retry: keep this conservative if you do retries in Push().
		MaxAttempts: 3,

		// Optional: surface kafka-go internal logs into your logging system.
		Logger:      kafka.LoggerFunc(func(msg string, a ...any) { log.Printf("[kafka-go] "+msg, a...) }),
		ErrorLogger: kafka.LoggerFunc(func(msg string, a ...any) { log.Printf("[kafka-go][err] "+msg, a...) }),
	}

	w := kafka.NewWriter(cfg)

	// Telecom-grade default: do NOT allow implicit topic creation.
	// Auto-create can create wrong partition counts/replication in prod.
	w.AllowAutoTopicCreation = false

	c.Writer = w

	log.Printf("kafka writer created: brokers=%v clientID=%s", brokers, clientID)

	// Topic lifecycle should be managed outside the producer (IaC / platform team).
	// If you MUST do it in-app, guard it behind a feature flag and ensure idempotency.
	if Configuration.Kafka_Events.CreateTopicsOnStartup {
		if err := ensureTopics(brokersCSVFromSlice(brokers)); err != nil {
			_ = c.Close()
			return err
		}
	}

	return nil
}

func ensureTopics(brokersCSV string) error {
	// Example only. Prefer IaC / separate init job in production.
	// Also: do NOT hardcode topic settings here; read from config.
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: Configuration.Kafka_Events.KafkaClientId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, topic := range LoyaltyKafkaTopics {
		err := CreateTopicIfMissing(
			ctx,
			brokersCSV,
			topic, // topic name
			12,    // partitions (production-grade example)
			Configuration.Kafka_Events.ReplicationFactor, // replication factor (HA recommended)
			"1209600000", // retention.ms (14 days)
			dialer,
			nil, // use default options
		)
		if err != nil {
			// Non-fatal: a broker/topic hiccup must not crash the loyalty service.
			// If brokers are truly down, per-message pushes will fail and be logged.
			log.Printf("failed to ensure %s topic exists (continuing): %v", topic, err)
			continue
		}
		log.Printf("topic %s ensured successfully", topic)
	}
	return nil
}

func splitAndTrimCSV(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		x := strings.TrimSpace(r)
		if x != "" {
			out = append(out, x)
		}
	}
	return out
}

func brokersCSVFromSlice(brokers []string) string {
	return strings.Join(brokers, ",")
}

// ***********************************************************************************************
// Kafka writer
// ***********************************************************************************************
type PushOptions struct {
	// Total time budget for a single Push call (all retries included).
	OverallTimeout time.Duration

	// Timeout per WriteMessages attempt (should be <= OverallTimeout).
	AttemptTimeout time.Duration

	// Max number of retries on transient errors (not counting the first attempt).
	MaxRetries int

	// Base backoff between retries; exponential backoff is applied.
	BaseBackoff time.Duration

	// Max backoff cap.
	MaxBackoff time.Duration

	// Optional: add a small random delay to avoid thundering herds.
	JitterFraction float64 // e.g. 0.2 means +/-20%

	// Optional: caller can provide a fixed timestamp (otherwise time.Now()).
	Now func() time.Time
}

func DefaultPushOptions() PushOptions {
	return PushOptions{
		OverallTimeout: 3 * time.Second,
		AttemptTimeout: 1 * time.Second,
		MaxRetries:     3,
		BaseBackoff:    50 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
		JitterFraction: 0.2,
		Now:            time.Now,
	}
}

func (c *KafkaClient) Push(ctx context.Context, key, value []byte, topic string, opt *PushOptions) error {
	if c == nil || c.Writer == nil {
		return errors.New("kafka writer is nil")
	}
	if topic == "" {
		return errors.New("topic is required")
	}

	options := DefaultPushOptions()
	if opt != nil {
		options = *opt
		if options.Now == nil {
			options.Now = time.Now
		}
	}

	// Budget the whole call, so it can't hang indefinitely.
	ctx, cancel := context.WithTimeout(ctx, options.OverallTimeout)
	defer cancel()

	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  options.Now(),
	}

	var lastErr error
	attempts := options.MaxRetries + 1

	for attempt := 0; attempt < attempts; attempt++ {
		// Respect caller cancellation / overall deadline.
		if err := ctx.Err(); err != nil {
			if lastErr != nil {
				return fmt.Errorf("push aborted (%v): last error: %w", err, lastErr)
			}
			return err
		}

		// Per-attempt timeout (bounded by remaining overall deadline automatically).
		attemptCtx, attemptCancel := context.WithTimeout(ctx, options.AttemptTimeout)
		err := c.Writer.WriteMessages(attemptCtx, msg)
		attemptCancel()

		if err == nil {
			return nil
		}
		lastErr = err

		// Don't retry non-transient failures.
		if !isRetryableKafkaGoError(err) || attempt == attempts-1 {
			return fmt.Errorf("kafka write failed (attempt %d/%d): %w", attempt+1, attempts, err)
		}

		// Backoff with jitter.
		sleep := backoffDuration(options.BaseBackoff, options.MaxBackoff, attempt)
		sleep = applyJitter(sleep, options.JitterFraction)

		// Sleep is also cancelable.
		timer := time.NewTimer(sleep)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("push canceled during backoff: %w", ctx.Err())
		case <-timer.C:
		}
	}

	return fmt.Errorf("kafka write failed after %d attempts: %w", attempts, lastErr)
}

func isRetryableKafkaGoError(err error) bool {
	// Retry common transient conditions.
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	// kafka-go exposes some sentinel errors.
	// LeaderNotAvailable often appears during elections / topic metadata refresh.
	if errors.Is(err, kafka.LeaderNotAvailable) {
		return true
	}

	// Many broker errors are wrapped; treat temporary network errors as retryable.
	type temporary interface{ Temporary() bool }
	var te temporary
	if errors.As(err, &te) && te.Temporary() {
		return true
	}

	// If you see more real-world transient errors in logs, add them here.
	// Example candidates (depending on kafka-go version/wrapping):
	// - kafka.RequestTimedOut
	// - io.EOF / net.Error (timeout)
	type netErr interface {
		Timeout() bool
		Temporary() bool
	}
	var ne netErr
	if errors.As(err, &ne) && (ne.Timeout() || ne.Temporary()) {
		return true
	}

	return false
}

func backoffDuration(base, max time.Duration, attempt int) time.Duration {
	// exp backoff: base * 2^attempt
	d := base << attempt
	if d > max {
		return max
	}
	return d
}

func applyJitter(d time.Duration, frac float64) time.Duration {
	if frac <= 0 {
		return d
	}
	// random in [1-frac, 1+frac]
	min := 1.0 - frac
	max := 1.0 + frac
	f := min + rand.Float64()*(max-min)
	return time.Duration(float64(d) * f)
}

// **************************************************************************************
// check and create topic
// **************************************************************************************
type TopicCreateOptions struct {
	// How long to wait for broker/controller operations.
	Timeout time.Duration

	// If true, we do a metadata check first.
	// Still idempotent (we also ignore "already exists" on create).
	CheckExistsFirst bool
}

func DefaultTopicCreateOptions() TopicCreateOptions {
	return TopicCreateOptions{
		Timeout:          10 * time.Second,
		CheckExistsFirst: true,
	}
}

// CreateTopicIfMissing creates a topic if it doesn't exist.
// It is safe to call concurrently from multiple instances (idempotent).
//
// retentionMs must be a string containing milliseconds, e.g. "1209600000" for 14 days.
func CreateTopicIfMissing(
	ctx context.Context,
	brokerURLsCSV string,
	topic string,
	partitions int,
	replication int,
	retentionMs string,
	dialer *kafka.Dialer,
	opt *TopicCreateOptions,
) error {
	if strings.TrimSpace(topic) == "" {
		return errors.New("topic is required")
	}
	if partitions <= 0 {
		return fmt.Errorf("invalid partitions=%d", partitions)
	}
	if replication <= 0 {
		return fmt.Errorf("invalid replication=%d", replication)
	}
	if strings.TrimSpace(brokerURLsCSV) == "" {
		return errors.New("brokers are not set")
	}

	// validate retentionMs is numeric
	if _, err := strconv.ParseInt(strings.TrimSpace(retentionMs), 10, 64); err != nil {
		return fmt.Errorf("retentionMs must be numeric milliseconds, got %q: %w", retentionMs, err)
	}

	options := DefaultTopicCreateOptions()
	if opt != nil {
		options = *opt
	}

	ctx, cancel := context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	brokers := splitAndTrimCSV(brokerURLsCSV)
	if len(brokers) == 0 {
		return errors.New("brokers list is empty after parsing")
	}

	if dialer == nil {
		dialer = &kafka.Dialer{Timeout: options.Timeout, DualStack: true}
	}

	// 1) Connect to ANY broker to find the controller
	var conn *kafka.Conn
	var err error
	for _, b := range brokers {
		conn, err = dialer.DialContext(ctx, "tcp", b)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("unable to dial any broker: %w", err)
	}
	defer conn.Close()

	// Optional existence check (nice for logging/avoids controller work),
	// but NOT sufficient alone due to races.
	if options.CheckExistsFirst {
		exists, exErr := topicExists(ctx, dialer, conn, topic)
		if exErr != nil {
			// In production, you usually log and continue to "create and ignore already exists"
			// rather than failing hard on a metadata hiccup.
			log.Printf("topicExists check failed for %s (continuing with create): %v", topic, exErr)
		} else if exists {
			return nil
		}
	}

	// 2) Get controller info
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	controllerConn, err := dialer.DialContext(ctx, "tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("failed to dial controller %s: %w", controllerAddr, err)
	}
	defer controllerConn.Close()

	// 3) Create topic (idempotent: ignore already exists)
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replication,
			ConfigEntries: []kafka.ConfigEntry{
				{ConfigName: "retention.ms", ConfigValue: strings.TrimSpace(retentionMs)},
			},
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		// IMPORTANT: if it already exists, treat it as success.
		// kafka-go typically returns kafka.TopicAlreadyExists (or wraps a protocol error).
		if errors.Is(err, kafka.TopicAlreadyExists) {
			return nil
		}
		return fmt.Errorf("create topic %s failed: %w", topic, err)
	}

	return nil
}

func topicExists(ctx context.Context, dialer *kafka.Dialer, conn *kafka.Conn, topic string) (bool, error) {
	// ReadPartitions hits cluster metadata and returns partitions for all topics.
	// Then we check if our topic appears.
	parts, err := conn.ReadPartitions()
	if err != nil {
		return false, err
	}
	for _, p := range parts {
		if p.Topic == topic {
			return true, nil
		}
	}
	return false, nil
}
