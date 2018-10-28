package fuzzy

import "math"

const(
	MustBeLowIncome = 500
	MustNotBeLowIncome = 1000

	MustNotBeMiddleLowIncome = 400
	MustBeMiddleLowIncome = 900
	MustBeMiddleHighIncome = 1200
	MustNotBeMiddleHighIncome = 1750

	MustNotBeHighIncome = 1100
	MustBeHighIncome = 1650
)

const(
	MustBeLowDebt = 15000
	MustNotBeLowDebt = 35000

	MustNotBeMiddleLowDebt = 10000
	MustBeMiddleLowDebt = 40000
	MustBeMiddleHighDebt = 50000
	MustNotBeMiddleHighDebt = 75000

	MustNotBeHighDebt = 45000
	MustBeHighDebt = 70000
)

const(
	AcceptedValue = 100
	ConsideredValue = 70
	RejectedValue = 50
)

type Fuzzy interface {
	Fuzzification(number *FuzzyNumber) error
	Defuzzification(number *FuzzyNumber) error
	Inference(number *FuzzyNumber) error
}

type Family struct {
	Number string
	Income float64
	Debt float64
}

type FuzzyNumber struct {
	Family Family

	IncomeMembership []float64
	DebtMembership []float64

	AccepetedInference float64
	ConsideredInference float64
	RejectedInference float64

	CrispValue float64
}

type BLT struct {
}

func (b *BLT) Inference(number *FuzzyNumber) error {
	number.AccepetedInference = math.Max(math.Min(number.IncomeMembership[0], number.DebtMembership[2]), math.Min(number.IncomeMembership[0], number.DebtMembership[1]))
	number.AccepetedInference = math.Max(number.AccepetedInference, math.Min(number.IncomeMembership[0], number.DebtMembership[0]))

	number.ConsideredInference = math.Max(math.Min(number.IncomeMembership[1], number.DebtMembership[1]), math.Min(number.IncomeMembership[1], number.DebtMembership[2]))
	number.ConsideredInference = math.Max(number.ConsideredInference, math.Min(number.IncomeMembership[2], number.DebtMembership[2]))

	number.RejectedInference =  math.Max(math.Min(number.IncomeMembership[2], number.DebtMembership[0]), math.Min(number.IncomeMembership[2], number.DebtMembership[1]))
	number.RejectedInference = math.Max(number.RejectedInference, math.Min(number.IncomeMembership[1], number.DebtMembership[0]))

	return nil
}

func (b *BLT) Defuzzification(number *FuzzyNumber) error {
	number.CrispValue = 0
	number.CrispValue += number.AccepetedInference*AcceptedValue
	number.CrispValue += number.ConsideredInference* ConsideredValue
	number.CrispValue += number.RejectedInference* RejectedValue
	number.CrispValue /= (number.AccepetedInference+number.ConsideredInference+number.RejectedInference)

	return nil
}

func (b *BLT) Fuzzification(number *FuzzyNumber) error {
	number.DebtMembership = append(number.DebtMembership, b.DebtLow(number.Family.Debt))
	number.DebtMembership = append(number.DebtMembership, b.DebtMiddle(number.Family.Debt))
	number.DebtMembership = append(number.DebtMembership, b.DebtHigh(number.Family.Debt))

	number.IncomeMembership = append(number.IncomeMembership, b.IncomeLow(number.Family.Income))
	number.IncomeMembership = append(number.IncomeMembership, b.IncomeMiddle(number.Family.Income))
	number.IncomeMembership = append(number.IncomeMembership, b.IncomeHigh(number.Family.Income))
	return nil
}

func (b *BLT) IncomeLow(income float64) float64 {
	if income <= MustBeLowIncome {
		return 1
	} else if income > MustNotBeLowIncome {
		return 0
	}
	return 1 - (float64(income - MustBeLowIncome) / float64(MustNotBeLowIncome - MustBeLowIncome))
}

func (b *BLT) IncomeMiddle(income float64) float64 {
	if income > MustBeMiddleLowIncome && income <= MustBeMiddleHighIncome {
		return 1
	} else if income < MustNotBeMiddleLowIncome || income > MustNotBeMiddleHighIncome {
		return 0
	} else if income < MustBeMiddleLowIncome && income >= MustNotBeMiddleLowIncome {
		return float64(income - MustNotBeMiddleLowIncome) / float64(MustBeMiddleLowIncome - MustNotBeMiddleLowIncome)
	}

	return 1 - float64(income - MustBeMiddleHighIncome) / float64(MustNotBeMiddleHighIncome - MustBeMiddleHighIncome)
}

func (b *BLT) IncomeHigh(income float64) float64 {
	if income <= MustNotBeHighIncome {
		return 0
	} else if income > MustBeHighIncome {
		return 1
	}
	return float64(income - MustNotBeHighIncome) / float64(MustBeHighIncome - MustNotBeHighIncome)
}


func (b *BLT) DebtLow(income float64) float64 {
	if income <= MustBeLowDebt {
		return 1
	} else if income > MustNotBeLowDebt {
		return 0
	}
	return 1-(float64(income - MustBeLowDebt) / float64(MustNotBeLowDebt - MustBeLowDebt))
}

func (b *BLT) DebtMiddle(income float64) float64 {
	if income > MustBeMiddleLowDebt && income <= MustBeMiddleHighDebt {
		return 1
	} else if income < MustNotBeMiddleLowDebt || income > MustNotBeMiddleHighDebt {
		return 0
	} else if income < MustBeMiddleLowDebt && income >= MustNotBeMiddleLowDebt {
		return float64(income - MustNotBeMiddleLowDebt) / float64(MustBeMiddleLowDebt - MustNotBeMiddleLowDebt)
	}

	return 1 - (float64(income - MustBeMiddleHighDebt) / float64(MustNotBeMiddleHighDebt - MustBeMiddleHighDebt))
}

func (b *BLT) DebtHigh(income float64) float64 {
	if income <= MustNotBeHighDebt {
		return 0
	} else if income > MustBeHighDebt {
		return 1
	}
	return float64(income - MustNotBeHighDebt) / float64(MustBeHighDebt - MustNotBeHighDebt)
}