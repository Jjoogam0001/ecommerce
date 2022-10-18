package model

import "github.com/jackc/pgtype"

type (
	Order struct {
		OrderNumber     int32       `json:"order_number"`
		OrderDate       pgtype.Date `json:"order_date"`
		RequiredDate    pgtype.Date `json:"required_date"`
		ShippedDate     pgtype.Date `json:"shipped_date"`
		Status          string      `json:"status"`
		Comments        string      `json:"comments"`
		Customer_number int32       `json:"customer_number"`
	}

	Customer struct {
		CustomerNumber    int    `json:"customer_number"`
		CustomerName      string `json:"customer_name"`
		ContactLastName   string `json:"contact_last_name"`
		ContactFirstName  string `json:"contact_first_name"`
		Phone             string `json:"phone"`
		AddressLine       string `json:"address_line"`
		AddressLine2      string `json:"address_line2"`
		City              string `json:"city"`
		State             string `json:"state"`
		PostalCode        string `json:"postal_code"`
		Country           string `json:"country"`
		SalesRepEmpNumber int32  `json:"sales_rep_employee_number"`
		CreditLimit       int64  `json:"credit_limit"`
	}

	Employee struct {
		EmployeeNumber int    `json:"employee_number"`
		LastName       string `json:"last_name"`
		FirstName      string `json:"first_name"`
		Extension      string `json:"extension"`
		Email          string `json:"email"`
		OfficeCode     string `json:"office_code"`
		ReportsTo      int    `json:"reports_to"`
		Job_Title      string `json:"job_title"`
	}

	Office struct {
		OfficeCode   string `json:"office_code"`
		City         string `json:"city"`
		Phone        string `json:"phone"`
		AddressLine  string `json:"address_line"`
		AddressLine2 string `json:"address_line2"`
		State        string `json:"state"`
		Country      string `json:"country"`
		PostalCode   string `json:"postal_code"`
		Territory    string `json:"territory"`
	}

	OrderDetail struct {
		OrderNumber     int32  `json:"order_number"`
		ProductCode     string `json:"product_code"`
		QuantityOrdered int32  `json:"quantity_ordered"`
		PriceEach       string `json:"price_each"`
		OrderLineNumber int8   `json:"order_line_number"`
	}
	Payment struct {
		CustomerNumber int32       `json:"customer_number"`
		CheckNumber    string      `json:"check_number"`
		PaymentDate    pgtype.Date `json:"payment_date"`
		Amount         float64     `json:"amount"`
	}
	Product struct {
		ProductCode        string  `json:"product_code"`
		ProductName        string  `json:"product_name"`
		ProductLine        string  `json:"product_line"`
		ProductScale       string  `json:"product_scale"`
		ProductVendor      string  `json:"product_vendor"`
		ProductDescription string  `json:"product_description"`
		QuantityInStock    int32   `json:"quantity_in_stock"`
		BuyPrice           float32 `json:"buy_price"`
		Msrp               float32 `json:"msrp"`
	}
	ProductLine struct {
		ProductLine     string `json:"product_line"`
		TextDescription string `json:"text_description"`
		HtmlDescription string `json:"html_description"`
		Image           []byte `json:"image"`
	}

	CustomerResponse struct {
		Customer Customer `json:"customer"`
		Status   string   `json:"status"`
	}
	ProductResponse struct {
		Product Product `json:"product"`
		Status  string  `json:"status"`
	}
	OrderDetailResponse struct {
		OrderDetail OrderDetail `json:"orderdetail"`
		Status      string      `json:"status"`
	}
	OfficeResponse struct {
		Office Office `json:"office"`
		Status string `json:"status"`
	}
	EmployeeResponse struct {
		Employee Employee `json:"employee"`
		Status   string   `json:"status"`
	}
	OrderResponse struct {
		Order  Order  `json:"order"`
		Status string `json:"status"`
	}
	PaymentResponse struct {
		Payment Payment `json:"payment"`
		Status  string  `json:"status"`
	}
)
