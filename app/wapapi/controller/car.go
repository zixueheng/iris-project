package controller

// Car 控制器
type Car struct {
}

// GetCarBy 根据 carname 查找用户
func (c *Car) GetCarBy(carname string) string {
	return "wapapi get: " + carname
}
