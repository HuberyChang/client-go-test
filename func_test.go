package main

/*
	在这个例子中，calculateFinalPrice函数能够根据传入的折扣策略（DiscountPredicate类型）动态地计算出不同用户等级的最终支付价格，
	展示了将方法定义为类型并作为参数使用的灵活性和优势
*/

// DiscountPredicate 定义一个折扣策略的函数类型
type DiscountPredicate func(userLevel string, totalPrice float64) float64

// 实现几个具体的折扣策略函数
func regularDiscount(userLevel string, totalPrice float64) float64 {
	if userLevel == "regular" {
		return totalPrice * 0.9 // 90% of total price for regular users
	}
	return totalPrice
}

func premiumDiscount(userLevel string, totalPrice float64) float64 {
	if userLevel == "premium" {
		return totalPrice * 0.8 // 80% of total price for premium users
	}
	return totalPrice
}

// 定义一个函数，它接受上述类型的函数作为参数
func calculateFinalPrice(discountStrategy DiscountPredicate, userLevel string, totalPrice float64) float64 {
	return discountStrategy(userLevel, totalPrice)
}

// 在main函数中，直接调用函数
func main() {
	totalPrice := 100.0

	// For a regular user
	finalPrice := calculateFinalPrice(regularDiscount, "regular", totalPrice)
	println("Regular User's Final Price: ", finalPrice)

	// For a premium user
	finalPrice = calculateFinalPrice(premiumDiscount, "premium", totalPrice)
	println("Premium User's Final Price: ", finalPrice)
}
