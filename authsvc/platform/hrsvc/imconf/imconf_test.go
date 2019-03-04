package imconf

/*
func TestSetConfig(t *testing.T) {

	var d map[string]map[string]map[string]interface{}

	d = make(map[string]map[string]map[string]interface{})
	var item1, item2 map[string]interface{}

	item1 = make(map[string]interface{})
	item2 = make(map[string]interface{})

	item1["database"] = "imdb"
	item1["host"] = "localhost"
	item1["port"] = 3306.0
	item1["user"] = "uniontest"
	item1["passwd"] = "uniontest"
	item2["database"] = "healthdb"
	item2["host"] = "localhost"
	item2["port"] = 3306.0
	item2["user"] = "uniontest"
	item2["passwd"] = "uniontest"

	var db map[string]map[string]interface{}

	db = make(map[string]map[string]interface{})
	db["healthdb"] = item2
	db["imdb"] = item1

	d["mysql"] = db

	var item3 map[string]interface{}

	item3 = make(map[string]interface{})

	item3["db"] = 0.0
	item3["host"] = "192.168.1.6"
	item3["port"] = 6379.0
	var queue map[string]map[string]interface{}

	queue = make(map[string]map[string]interface{})
	queue["queue"] = item3
	d["redis"] = queue

	var item4, item5, item6 map[string]interface{}

	item4 = make(map[string]interface{})
	item5 = make(map[string]interface{})
	item6 = make(map[string]interface{})

	item4["host"] = "localhost"
	item4["port"] = 12000.0
	var service map[string]map[string]interface{}

	service = make(map[string]map[string]interface{})
	service["tcpservice"] = item4

	item5["host"] = "localhost"
	item5["port"] = 12001.0
	service["udpservice"] = item5

	item6["host"] = "localhost"
	item6["port"] = 7890.0
	service["httpservice"] = item6

	d["service"] = service

	var config ImConf
	

	buf, _ := json.Marshal(d)

	config.SetConfig(d)

	fmt.Printf("the configuration is: %v\n", config)

	fmt.Printf("the marshal map: %s\n", buf)

}

*/
