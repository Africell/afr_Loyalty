package Lendme

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"redisx"
)

// batchSpec describes one points-detail batch to create.
type batchSpec struct {
	YearMonth string // "YYYYMM"
	Awarded   float64
	Redeemed  float64
}

// accountSpec describes the seeded state for one test account.
type accountSpec struct {
	MSISDN     string
	Level      string    // Loyalty_Level_Key (all test plans live under SILVER)
	Segment    string    // Loyalty_Account_Segment_Key -> plan -> expiry rule
	OptStatus  string    // "OptedIn" | "OptedOut"
	LastOpt    time.Time // Last_Opt_Status_Date / Opt_Status_Date (zero = leave for OptedIn)
	FirstOptIn time.Time // First_Opt_In_Status_Date (refDate for legacy/fixed rules)
	Batches    []batchSpec
	Note       string // expected outcome, for the printout
}

func d(y, m, day int) time.Time {
	return time.Date(y, time.Month(m), day, 0, 0, 0, 0, time.UTC)
}

// Common dates.
var (
	firstOptIn = d(2024, 9, 15) // legacy refDate -> Initial 2025-09-15, Coming 2025-12-15
)

// Per-rule batch sets (one batch that expires now + one that is kept).
var (
	monthlyBatches   = []batchSpec{{"202412", 20, 0}, {"202601", 30, 0}} // exp 2025-06-30 / keep 2026-07-31
	quarterlyBatches = []batchSpec{{"202412", 20, 0}, {"202601", 30, 0}} // exp 2025-12-31 / keep 2027-03-31
	// Legacy/yearly: Initial = FirstOptIn + 12mo = 2025-09-15 -> expire Year_Month <= 2025-09 (inclusive).
	legacyBatches = []batchSpec{
		{"202401", 10, 0}, // Jan 2024  <= 2025-09 -> EXPIRE
		{"202412", 20, 0}, // Dec 2024  <= 2025-09 -> EXPIRE
		{"202509", 22, 0}, // Sep 2025  == Initial month -> EXPIRE (inclusive boundary)
		{"202510", 15, 0}, // Oct 2025  >  2025-09 -> KEEP (first keeper)
		{"202601", 25, 0}, // Jan 2026  -> KEEP
		{"202603", 18, 0}, // Mar 2026  -> KEEP
		{"202610", 12, 0}, // Oct 2026  -> KEEP (beyond Initial 2025-09; not in this cycle)
	} // OptedIn/legacy: 52 expire (202401+202412+202509) / 70 keep
)

// uatSeedMatrix builds the 16 account specs (8 segments x opted-in/opted-out) for the
// given fixed opt-out date. With optOut = 2026-04-25 the Fixed opted-out accounts resolve
// a Coming_Expiry of 2026-07-24 (opt-out + 90 days) / 2026-07-25 (opt-out + 3 months);
// once those dates pass, all of the account's points expire rather than being held.
func uatSeedMatrix(optOut time.Time) []accountSpec {
	return []accountSpec{
		// ── Monthly (Validity 6mo) ──
		{MSISDN: "2201700005", Level: "SILVER", Segment: "TEST_MONTHLY_FOLLOW", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: monthlyBatches, Note: "Monthly 6mo, OptedIn: 202412 expire / 202601 keep"},
		{MSISDN: "2201700001", Level: "SILVER", Segment: "TEST_MONTHLY_FOLLOW", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: monthlyBatches, Note: "Monthly+FollowOptedIn, recent opt-out (Phase3 skips): follows monthly -> 202412 expire"},

		{MSISDN: "2201000000", Level: "SILVER", Segment: "TEST_MONTHLY_FIXED_DAY", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: monthlyBatches, Note: "Monthly 6mo, OptedIn: 202412 expire / 202601 keep"},
		{MSISDN: "2201000101", Level: "SILVER", Segment: "TEST_MONTHLY_FIXED_DAY", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: monthlyBatches, Note: "OptedOut+Fixed 90d -> Coming = opt-out + 90d = 2026-07-24: all 50 expire once past"},

		{MSISDN: "2201700002", Level: "SILVER", Segment: "TEST_MONTHLY_FIXED_MONTH", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: monthlyBatches, Note: "Monthly 6mo, OptedIn: 202412 expire / 202601 keep"},
		{MSISDN: "2200090392", Level: "SILVER", Segment: "TEST_MONTHLY_FIXED_MONTH", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: monthlyBatches, Note: "OptedOut+Fixed 3mo -> Coming = opt-out + 3mo = 2026-07-25: all 50 expire once past"},

		// ── Quarterly ──
		{MSISDN: "2201700006", Level: "SILVER", Segment: "TEST_QUARTERLY_FIXED_DAY", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: quarterlyBatches, Note: "Quarterly, OptedIn: 202412 expire (2025-12-31) / 202601 keep (2027-03-31)"},
		{MSISDN: "2201700004", Level: "SILVER", Segment: "TEST_QUARTERLY_FIXED_DAY", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: quarterlyBatches, Note: "OptedOut+Fixed 90d -> Coming = opt-out + 90d = 2026-07-24: all 50 expire once past"},

		{MSISDN: "2201700000", Level: "SILVER", Segment: "TEST_QUARTERLY_FIXED_MONTH", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: quarterlyBatches, Note: "Quarterly, OptedIn: 202412 expire / 202601 keep"},
		{MSISDN: "2201700007", Level: "SILVER", Segment: "TEST_QUARTERLY_FIXED_MONTH", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: quarterlyBatches, Note: "OptedOut+Fixed 3mo -> Coming = opt-out + 3mo = 2026-07-25: all 50 expire once past"},

		// ── Yearly / legacy (Opted_In_Rule_Type == "", Validity 12mo + Grace 3mo) ──
		{MSISDN: "2201700016", Level: "SILVER", Segment: "TEST_YEARLY_FIXED_DAY", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: legacyBatches, Note: "Legacy 12mo+3, Initial=2025-09: 52 expire (<=202509) / 70 keep, Coming 2025-12-15"},
		{MSISDN: "2201700011", Level: "SILVER", Segment: "TEST_YEARLY_FIXED_DAY", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: legacyBatches, Note: "OptedOut+Fixed (Day window) -> Coming = opt-out + 90d = 2026-07-24: all 122 expire once past"},

		{MSISDN: "2201700012", Level: "SILVER", Segment: "TEST_YEARLY_FIXED_MONTH", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: legacyBatches, Note: "Legacy 12mo+3, Initial=2025-09: 52 expire (<=202509) / 70 keep"},
		{MSISDN: "2201700019", Level: "SILVER", Segment: "TEST_YEARLY_FIXED_MONTH", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: legacyBatches, Note: "OptedOut+Fixed (Month window) -> Coming = opt-out + 3mo = 2026-07-25: all 122 expire once past"},

		{MSISDN: "2201700010", Level: "SILVER", Segment: "TEST_YEARLY_FOLLOW", OptStatus: "OptedIn", FirstOptIn: firstOptIn,
			Batches: legacyBatches, Note: "Legacy 12mo+3, Initial=2025-09: 52 expire (<=202509) / 70 keep"},
		{MSISDN: "2201700009", Level: "SILVER", Segment: "TEST_YEARLY_FOLLOW", OptStatus: "OptedOut", LastOpt: optOut, FirstOptIn: firstOptIn,
			Batches: legacyBatches, Note: "Legacy+FollowOptedIn, recent opt-out (Phase3 skips): follows legacy -> 52 expire / 70 keep"},
	}
}

func seedAccount(ctx context.Context, spec accountSpec) error {
	// 1. read existing account, or build a fresh one (collections may be empty)
	account, err := getJSONWithMongoFallback[Customer_Loyalty_Account](ctx, Customer_Loyalty_Account{Key: spec.MSISDN}.RedisKey(), Mdb_Customer_Loyalty_Account, bson.M{"Key": spec.MSISDN})
	if err != nil {
		if !redisx.IsNil(err) {
			return fmt.Errorf("read account: %w", err)
		}
		cid, _ := strconv.ParseInt(spec.MSISDN, 10, 64)
		account = Customer_Loyalty_Account{
			Key:            spec.MSISDN,
			Customer_Id:    cid,
			Creation_date:  d(2024, 1, 1),
			Account_Status: "active",
		}
	}

	// 2. assign level + segment (this is what maps the account to its expiry rule)
	now := time.Now().UTC()
	account.Loyalty_Level_Key = spec.Level
	account.Loyalty_Level_Date = now
	account.Loyalty_Account_Segment_Key = spec.Segment
	account.Loyalty_Account_Segment_Date = now

	// 3. resolve the rule via plan -> expiry rule (tolerant: a missing plan/rule
	//    is itself a seeded edge case — the live process will log the error path)
	planKey := account.Loyalty_Level_Key + "|" + account.Loyalty_Account_Segment_Key
	var rule Loyalty_Point_Expiry_Rules
	var ruleKey string
	planMissing := false
	plan, planErr := redisx.GetJSON[Loyalty_Plan](ctx, RedisClient, Loyalty_Plan{Key: planKey}.RedisKey())
	if planErr != nil {
		planMissing = true
	} else {
		ruleKey = plan.Expiry_Rules_Key
		r, ruleErr := redisx.GetJSON[Loyalty_Point_Expiry_Rules](ctx, RedisClient, Loyalty_Point_Expiry_Rules{Key: ruleKey}.RedisKey())
		if ruleErr != nil {
			planMissing = true
		} else {
			rule = r
		}
	}

	// 4. delete any pre-existing batches for this MSISDN (clean slate)
	if _, err := Mdb_Customer_Loyalty_Account_Points_Detail.Coll.DeleteMany(ctx, bson.M{"Key": bson.M{"$regex": "^" + spec.MSISDN + "\\|"}}); err != nil {
		return fmt.Errorf("delete old batches (mongo): %w", err)
	}
	if _, err := redisx.DeleteByPattern(ctx, RedisClient, "Customer_Loyalty_Account_Points_Detail:"+spec.MSISDN+"|*", 200, true); err != nil {
		log.Printf("MSISDN %s: pattern delete warning: %v", spec.MSISDN, err)
	}

	// 5. build the new batches
	var batches []Customer_Loyalty_Account_Points_Detail
	var sumAwarded, sumRedeemed float64
	keys := []string{}
	for _, b := range spec.Batches {
		bd := Customer_Loyalty_Account_Points_Detail{
			Key:              spec.MSISDN + "|" + b.YearMonth,
			Year_Month:       b.YearMonth,
			Creation_date:    monthStart(b.YearMonth),
			Awarded_Points:   b.Awarded,
			Redeemed_Points:  b.Redeemed,
			Available_Points: b.Awarded - b.Redeemed,
			Expired_Points:   0,
			Last_Credit_Date: monthStart(b.YearMonth),
		}
		batches = append(batches, bd)
		sumAwarded += b.Awarded
		sumRedeemed += b.Redeemed
		keys = append(keys, bd.Key)
	}

	// 6. set account state per the spec
	account.Opt_Status = spec.OptStatus
	if !spec.LastOpt.IsZero() {
		// Last_Opt_Status_Date maps to bson "Opt_Status_Date", which Phase 3 filters on.
		account.Last_Opt_Status_Date = spec.LastOpt
	}
	if !spec.FirstOptIn.IsZero() {
		account.First_Opt_In_Status_Date = spec.FirstOptIn
	}
	account.Expiry_Date = time.Time{} // force fresh refDate (no drift) on resolve/legacy
	account.Awarded_Points = sumAwarded
	account.Redeemed_Points = sumRedeemed
	account.Available_Points = sumAwarded - sumRedeemed
	account.Expired_Points = 0
	account.Redeemed_Expired_Points = 0
	account.Outstanding_fraction_points = 0
	account.Points_Detail_Keys = keys
	account.Last_Award_Date = now

	// 7. compute Coming_Expiry_Date / Initial_Date exactly like the live app does.
	//    - misconfigured (no plan): force a past Coming so the scheduler picks it up.
	//    - Monthly/Quarterly/OptedOut+Fixed: resolveAccountExpiryDates.
	//    - legacy (empty Opted_In_Rule_Type): mirror the Customer_Loyalty_Account_Get
	//      legacy branch (resolve returns zero for it).
	var comingExpiry, initialDate time.Time
	var pointsToExpire float64
	if planMissing {
		comingExpiry = d(2025, 1, 1)
	} else {
		comingExpiry, initialDate, pointsToExpire = resolveAccountExpiryDates(account, rule, batches)
		if comingExpiry.IsZero() {
			// legacy fixed-date path
			refDate := account.Expiry_Date
			if refDate.IsZero() {
				refDate = account.First_Opt_In_Status_Date
			}
			if refDate.IsZero() {
				refDate = account.Creation_date
			}
			initialDate = addValidity(refDate, rule.Validity_Unit, rule.Validity_Duration)
			comingExpiry = addValidity(initialDate, rule.Grace_Validity_Unit, rule.Grace_Validity_Duration)
			pointsToExpire = 0
			for _, b := range batches {
				if len(b.Year_Month) != 6 {
					continue
				}
				y, yErr := strconv.Atoi(b.Year_Month[:4])
				m, mErr := strconv.Atoi(b.Year_Month[4:])
				if yErr != nil || mErr != nil {
					continue
				}
				if y < initialDate.Year() || (y == initialDate.Year() && m <= int(initialDate.Month())) {
					pointsToExpire += b.Available_Points
				}
			}
		}
	}
	account.Coming_Expiry_Date = comingExpiry
	account.Initial_Date = initialDate
	account.Points_To_Expire = pointsToExpire

	// 8. write batches (Mongo + Redis)
	for _, bd := range batches {
		if _, err := Mdb_Customer_Loyalty_Account_Points_Detail.Coll.UpdateOne(ctx,
			bson.M{"Key": bd.Key}, bson.M{"$set": bd}, options.UpdateOne().SetUpsert(true)); err != nil {
			return fmt.Errorf("upsert batch %s (mongo): %w", bd.Key, err)
		}
		if err := redisx.SetJSON(ctx, RedisClient, bd.RedisKey(), bd); err != nil {
			return fmt.Errorf("set batch %s (redis): %w", bd.Key, err)
		}
	}

	// 9. write account (Mongo + Redis)
	if _, err := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(ctx,
		bson.M{"Key": account.Key}, bson.M{"$set": account}, options.UpdateOne().SetUpsert(true)); err != nil {
		return fmt.Errorf("upsert account (mongo): %w", err)
	}
	if err := redisx.SetJSON(ctx, RedisClient, account.RedisKey(), account); err != nil {
		return fmt.Errorf("set account (redis): %w", err)
	}

	// 10. Reflect the seeded points in governance exactly like the live award flow:
	//     Loyalty_Governance_Available_Points_Debit increments Distributed_Points_Pool
	//     by the awarded points. The spendable balance is derived
	//     (Available_Points_Pool - Distributed - Expired - Redeemed), so this removes
	//     the seeded points from what's available and books them as distributed,
	//     without mutating the Available_Points_Pool cap. Non-fatal if governance is missing.
	if sumAwarded > 0 {
		var uc UserControl
		if gErr := uc.Loyalty_Governance_Available_Points_Debit(sumAwarded); gErr != nil {
			log.Printf("MSISDN %s: governance distribute warning: %v", spec.MSISDN, gErr)
		}
	}

	ruleLabel := ruleKey
	if planMissing {
		ruleLabel = "<no plan: " + planKey + ">"
	}
	log.Printf("MSISDN %-3s | %-8s | seg=%-26s | rule=%-28s | avail=%-5.0f | Initial=%s | Coming=%s | toExpire=%.0f\n        -> %s",
		spec.MSISDN, account.Opt_Status, spec.Segment, ruleLabel, account.Available_Points,
		fmtDate(initialDate), fmtDate(comingExpiry), pointsToExpire, spec.Note)
	return nil
}

func monthStart(yearMonth string) time.Time {
	t, err := time.Parse("200601", yearMonth)
	if err != nil {
		return time.Time{}
	}
	return t.UTC()
}

func fmtDate(t time.Time) string {
	if t.IsZero() {
		return "<zero>"
	}
	return t.Format("2006-01-02")
}

// SeedUATLoyaltyTestPointsExpiry seeds the points-expiry edge cases on UAT:
// resets governance pools, creates the 8 TEST_* segments/expiry rules/plans, then
// re-seeds 16 fixed accounts (both Mongo + Redis) so the live expiry job can be
// observed without a restart. Assumes Mongo repos + RedisClient are already
// initialized (true in main). Best-effort throughout; errors are logged, not fatal.
func (Uc *UserControl) SeedUATLoyaltyTestPointsExpiry() {
	if Configuration.IsLoyaltyProduction {
		return
	}

	log.Println("<<SeedUATLoyaltyTestPointsExpiry>> started")
	Uc.resetGovernancePools()
	Uc.seedUATTestConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// Fixed opt-out date (not derived from time.Now) so the seeded opt-out date — and
	// therefore the expected outcomes documented for the commercial team — stay identical
	// on every deployment.
	for _, spec := range uatSeedMatrix(d(2026, 4, 25)) {
		if err := seedAccount(ctx, spec); err != nil {
			log.Printf("<<SeedUATLoyaltyTestPointsExpiry>> MSISDN %s: %v", spec.MSISDN, err)
		}
	}
	log.Println("<<SeedUATLoyaltyTestPointsExpiry>> finished")
}

const seedLogin = "SeedUATExpiry"

// uatTestPlanSpec ties a test segment to the expiry rule its plan (SILVER|<Segment>) points at.
type uatTestPlanSpec struct {
	Segment string
	Rule    Loyalty_Point_Expiry_Rules_AddRequest
}

// The 8 expiry-rule edge cases: Monthly(6mo)/Quarterly/Yearly(12mo+3 grace),
// each crossed with opt-out behavior FollowOptedIn / Fixed Day(90) / Fixed Month(3).
var uatTestPlanSpecs = []uatTestPlanSpec{
	{Segment: "TEST_MONTHLY_FOLLOW", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_MONTHLY_FOLLOW", Description: "TEST Monthly 6mo, opt-out follows opted-in",
		Opted_In_Rule_Type: "Monthly", Validity_Unit: "Month", Validity_Duration: 6,
		Opted_Out_Rule_Type: "FollowOptedIn"}},
	{Segment: "TEST_MONTHLY_FIXED_DAY", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_MONTHLY_FIXED_DAY", Description: "TEST Monthly 6mo, opt-out fixed 90 days",
		Opted_In_Rule_Type: "Monthly", Validity_Unit: "Month", Validity_Duration: 6,
		Opted_Out_Rule_Type: "Fixed", Opted_Out_Validity_Unit: "Day", Opted_Out_Validity_Duration: 90}},
	{Segment: "TEST_MONTHLY_FIXED_MONTH", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_MONTHLY_FIXED_MONTH", Description: "TEST Monthly 6mo, opt-out fixed 3 months",
		Opted_In_Rule_Type: "Monthly", Validity_Unit: "Month", Validity_Duration: 6,
		Opted_Out_Rule_Type: "Fixed", Opted_Out_Validity_Unit: "Month", Opted_Out_Validity_Duration: 3}},
	{Segment: "TEST_QUARTERLY_FIXED_DAY", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_QUARTERLY_FIXED_DAY", Description: "TEST Quarterly, opt-out fixed 90 days",
		Opted_In_Rule_Type: "Quarterly",
		Opted_Out_Rule_Type: "Fixed", Opted_Out_Validity_Unit: "Day", Opted_Out_Validity_Duration: 90}},
	{Segment: "TEST_QUARTERLY_FIXED_MONTH", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_QUARTERLY_FIXED_MONTH", Description: "TEST Quarterly, opt-out fixed 3 months",
		Opted_In_Rule_Type: "Quarterly",
		Opted_Out_Rule_Type: "Fixed", Opted_Out_Validity_Unit: "Month", Opted_Out_Validity_Duration: 3}},
	{Segment: "TEST_YEARLY_FIXED_DAY", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_YEARLY_FIXED_DAY", Description: "TEST Yearly legacy 12mo+3, opt-out fixed 90 days",
		Opted_In_Rule_Type: "", Validity_Unit: "Month", Validity_Duration: 12,
		Grace_Validity_Unit: "Month", Grace_Validity_Duration: 3,
		Opted_Out_Rule_Type: "Fixed", Opted_Out_Validity_Unit: "Day", Opted_Out_Validity_Duration: 90}},
	{Segment: "TEST_YEARLY_FIXED_MONTH", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_YEARLY_FIXED_MONTH", Description: "TEST Yearly legacy 12mo+3, opt-out fixed 3 months",
		Opted_In_Rule_Type: "", Validity_Unit: "Month", Validity_Duration: 12,
		Grace_Validity_Unit: "Month", Grace_Validity_Duration: 3,
		Opted_Out_Rule_Type: "Fixed", Opted_Out_Validity_Unit: "Month", Opted_Out_Validity_Duration: 3}},
	{Segment: "TEST_YEARLY_FOLLOW", Rule: Loyalty_Point_Expiry_Rules_AddRequest{
		Key: "TEST_EXP_YEARLY_FOLLOW", Description: "TEST Yearly legacy 12mo+3, opt-out follows opted-in",
		Opted_In_Rule_Type: "", Validity_Unit: "Month", Validity_Duration: 12,
		Grace_Validity_Unit: "Month", Grace_Validity_Duration: 3,
		Opted_Out_Rule_Type: "FollowOptedIn"}},
}

// seedUATTestConfig creates the 8 test segments, expiry rules, and SILVER|<segment>
// plans. Best-effort: existing keys no-op. The SILVER level entity is intentionally
// not created (Loyalty_Level_Add rejects intersecting ranges and the expiry read
// path only needs the plan).
func (Uc *UserControl) seedUATTestConfig() {
	for _, s := range uatTestPlanSpecs {
		Uc.Loyalty_Account_Segment_Add(seedLogin, Loyalty_Account_Segment_AddRequest{
			Key: s.Segment, Description: s.Segment,
			Amount_From: 0, Amount_Till: 999999999, AON_From: 0, AON_Till: 999999999,
		})
		Uc.Loyalty_Point_Expiry_Rules_Add(seedLogin, s.Rule)
		Uc.Loyalty_Plan_Add(seedLogin, Loyalty_Plan_AddRequest{
			Key:                         "SILVER|" + s.Segment,
			Description:                 s.Segment + " test plan",
			Loyalty_Level_Key:           "SILVER",
			Loyalty_Account_Segment_Key: s.Segment,
			Earning_Rules_Key:           "Default_Earning_Rules",
			Expiry_Rules_Key:            s.Rule.Key,
			Redemption_Rules_Key:        "Default_Redemption_Rules",
		})
	}
}

// resetGovernancePools zeroes the distributed/redeemed/expired pools so the points
// re-added by seeding are idempotent. Preserves Available_Points_Pool when the entry
// exists; creates a default otherwise.
func (Uc *UserControl) resetGovernancePools() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	gov, gErr := redisx.GetJSON[Loyalty_Governance](ctx, RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
	if gErr != nil {
		log.Printf("resetGovernancePools: governance missing/unreadable (%v) — creating default", gErr)
		gov = Loyalty_Governance{
			Key:                             LOYALTY_GOVERNANCE_KEY,
			Available_Points_Pool:           5000000000,
			MaxAllowedPoints_PerTransaction: 100,
			MaxSubsAwardedPoints_PerMonth:   10000,
			MaxSubsAwardedPoints:            100000,
			DailyEarningLimit:               10000,
			DailyPointsRedemptionLimit:      10000,
			DailyRedemptionAttemptLimit:     10,
			WeeklyEarningLimit:              100000,
			WeeklyPointsRedemptionLimit:     100000,
			WeeklyRedemptionAttemptLimit:    100,
		}
	}
	gov.Distributed_Points_Pool = 0
	gov.Redeemed_Points_Pool = 0
	gov.Expired_Points_Pool = 0
	if _, uErr := Mdb_Loyalty_Governance.Coll.UpdateOne(ctx, bson.M{"Key": gov.Key}, bson.M{"$set": gov}, options.UpdateOne().SetUpsert(true)); uErr != nil {
		log.Printf("resetGovernancePools: mongo warning: %v", uErr)
	}
	if sErr := redisx.SetJSON(ctx, RedisClient, gov.RedisKey(), gov); sErr != nil {
		log.Printf("resetGovernancePools: redis warning: %v", sErr)
	}
}
