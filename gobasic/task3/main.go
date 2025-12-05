package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wangxg-sj/web3-learning/gobasic/task3/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//curd()
	//transaction(1, 2, 500)
	//sqlX1()
	//sqlX2()
	//gorm1()
	gorm2()
}

type Student struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:20"`
	Age   uint8
	Grade string `gorm:"size:20"`
}

// 1
func curd() {
	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/web3?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	// 自动迁移
	err = db.AutoMigrate(&Student{})
	if err != nil {
		fmt.Printf("数据库迁移失败: %v", err)
	}
	student1 := Student{
		Name:  "张三",
		Age:   18,
		Grade: "三年级",
	}
	ctx := context.Background()
	tx := db.WithContext(ctx).Create(&student1)
	if tx.Error != nil {
		fmt.Printf("创建失败: %v", tx.Error)
	}
	fmt.Printf("创建成功: %v 条数据 \n", tx.RowsAffected)

	find, err := gorm.G[Student](db).Where("age >= ?", 18).Find(ctx)
	if err != nil {
		fmt.Printf("查询失败: %v", err)
	}
	for _, row := range find {
		fmt.Println(row)
	}

	//
	affected, err := gorm.G[Student](db).Where("name = ?", "张三").Updates(ctx, Student{
		Age: 20,
	})
	if err != nil {
		fmt.Printf("更新失败: %v\n", err)
	}
	fmt.Printf("更新成功: %v 条数据 \n", affected)

	t := db.Model(&Student{}).Where("age < ?", 21).Delete(&Student{})
	if t.Error != nil {
		fmt.Printf("删除失败: %v\n", t.Error)
	}
	fmt.Printf("删除成功: %v 条数据 \n", t.RowsAffected)
}

func transaction(fromId uint, toId uint, amount float64) {
	type Account struct {
		gorm.Model
		ID      uint `gorm:"primaryKey;autoIncrement:false"`
		Balance float64
	}
	type Transaction struct {
		gorm.Model
		ID     uint `gorm:"primaryKey"`
		From   uint
		To     uint
		Amount float64
	}

	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/web3?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	// 自动迁移
	//err = db.AutoMigrate(&Account{}, &Transaction{})
	//// 创建两个账号
	//account1 := Account{
	//	ID:      1,
	//	Balance: 1000,
	//}
	//account2 := Account{
	//	ID:      2,
	//	Balance: 2000,
	//}
	background := context.Background()
	//db.WithContext(background).Create(&account1)
	//db.WithContext(background).Create(&account2)
	err = db.Transaction(func(tx *gorm.DB) error {
		//检查account1余额
		var balance float64
		tx.WithContext(background).Raw("select balance from web3.accounts where id = ?", fromId).
			Scan(&balance)
		if balance < amount {
			return fmt.Errorf("account %d balance not enough", fromId)
		}
		// 检查toId是否存在
		var count int64
		tx.WithContext(background).Raw("select count(*) from web3.accounts where id = ?", toId).
			Scan(&count)
		if count == 0 {
			return fmt.Errorf("account %d not found", toId)
		}
		_, err1 := tx.WithContext(background).
			Raw("update accounts set balance = balance - ? where id = ?", amount, fromId).Rows()

		if err1 != nil {
			return err1
		}

		_, err2 := tx.WithContext(background).
			Raw("update accounts set balance = balance + ? where id = ?", amount, toId).Rows()

		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		fmt.Printf("事务失败: %v\n", err)
	}
}

func sqlX1() {
	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/web3?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	type Employee struct {
		ID         int     `db:"id"`
		Name       string  `db:"name"`
		Department string  `db:"department"`
		Salary     float64 `db:"salary"`
	}
	var employees []Employee
	// 背景上下文
	background := context.Background()
	err = db.SelectContext(background, &employees, "select * from employees where department = ?", "技术部")
	if err != nil {
		fmt.Println("查询失败:", err)
	}
	for _, emp := range employees {
		fmt.Println(emp)
	}

	var employee Employee
	// 背景上下文
	err = db.GetContext(background, &employee, "select * from employees order by salary desc limit 1")
	if err != nil {
		fmt.Println("查询失败:", err)
	}
	fmt.Printf("最高工资员工: %v\n", employee)
}

func sqlX2() {
	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/web3?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	type Book struct {
		ID     int     `db:"id"`
		Title  string  `db:"title"`
		Author string  `db:"author"`
		Price  float64 `db:"price"`
	}
	_, err = db.ExecContext(context.Background(),
		"create table if not exists books (id int primary key auto_increment, title varchar(255), author varchar(255), price float)")
	if err != nil {
		fmt.Println("创建表失败:", err)
		return
	}

	// 批量插入
	_, err = db.ExecContext(context.Background(),
		"insert into books (title, author, price) values (?, ?, ?), (?, ?, ?), (?, ?, ?)",
		"Go 语言圣经", "Rob Pike", 50,
		"Go 语言高级编程", "柴树杉", 60,
		"Go 语言设计模式", "王方", 70)
	if err != nil {
		fmt.Println("插入失败:", err)
		return
	}

	var books []Book
	background := context.Background()
	err = db.SelectContext(background, &books,
		"select * from books where price >= ? ", 50)
	if err != nil {
		fmt.Println("查询失败:", err)
	}
	for _, book := range books {
		fmt.Println(book)
	}

}

func gorm1() {
	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/web3?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Comment{}, &model.Post{})
	if err != nil {
		fmt.Println("数据库迁移失败:", err)
		return
	}

	//创建三个用户并分别创建两篇文章,第一篇文章有两条评论,第二篇文章有一条评论

	users := []model.User{
		{Name: "张三", Age: 30},
		{Name: "李四", Age: 25},
		{Name: "王五", Age: 35},
	}
	db.Create(&users)
	//创建两篇文章
	posts := []model.Post{
		{Title: "第一篇文章", Author: "张三", Content: "这是第一篇文章的内容", UserId: 1},
		{Title: "第二篇文章", Author: "李四", Content: "这是第二篇文章的内容", UserId: 2},
	}
	db.Create(&posts)
	//创建评论
	comments := []model.Comment{
		{PostID: 1, Author: "李四", Content: "这是第一篇文章的第一条评论"},
		{PostID: 1, Author: "王五", Content: "这是第一篇文章的第二条评论"},
		{PostID: 2, Author: "张三", Content: "这是第二篇文章的第一条评论"},
	}
	db.Create(&comments)

}

func gorm2() {
	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/web3?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}

	//查询用户的所有文章和评论
	var user model.User
	db.Preload("Posts").Preload("Posts.Comments").First(&user, 1)
	//json序列化 格式化输出
	jsonStr, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		fmt.Println("json序列化失败:", err)
	}
	fmt.Println(string(jsonStr))

	// 查询评论最多的文章

	type PostWithCommentCount struct {
		model.Post
		CommentCount int `gorm:"column:comment_count"`
	}
	var post PostWithCommentCount
	db.Debug().Table("posts as p").
		Select("p.*, count(c.id) as comment_count").
		Joins("left join comments as c on p.id = c.post_id").
		Group("p.id").
		Order("comment_count desc").
		First(&post)

	fmt.Printf("评论最多的文章:%v, 评论数:%d\n", post.ID, post.CommentCount)

	//删除第一篇文章的第一条评论
	var comment model.Comment
	db.First(&comment, 3)
	db.Delete(&comment)

}
