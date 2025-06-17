package tui

type OrderDetail [][]string

func GetOrderDetail() *OrderDetail {
	var orderDetail OrderDetail

	orderDetail = append(orderDetail, []string{"charm", "200"})
	orderDetail = append(orderDetail, []string{"comics", "250"})
	orderDetail = append(orderDetail, []string{"failed notes", "100"})

	return &orderDetail
}
