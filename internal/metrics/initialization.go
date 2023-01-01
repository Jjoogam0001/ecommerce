package metrics

// API metrics counters

type APIRequestCountLabels struct {
	endpoint string
	method   string
}

func initializeApiMetricsCounters() {
	originators := []string{}
	apiRequestCountLabels := []APIRequestCountLabels{
		{
			endpoint: "/customers/delete",
			method:   methodDelete,
		},
		{
			endpoint: "/customers/get",
			method:   methodGet,
		},
		{
			endpoint: "/customers/customers",
			method:   methodGet,
		},
		{
			endpoint: "/customers/update",
			method:   methodPatch,
		},
		{
			endpoint: "/employees/delete",
			method:   methodDelete,
		},
		{
			endpoint: "/employees/get",
			method:   methodGet,
		},
		{
			endpoint: "/employees/employees",
			method:   methodGet,
		},
		{
			endpoint: "/employees/update",
			method:   methodPatch,
		},
		{
			endpoint: "/offices/delete",
			method:   methodDelete,
		},
		{
			endpoint: "/offices/get",
			method:   methodGet,
		},
		{
			endpoint: "/offices/offices",
			method:   methodGet,
		},
		{
			endpoint: "/offices/update",
			method:   methodPatch,
		},
		{
			endpoint: "/orderDetails/delete",
			method:   methodDelete,
		},
		{
			endpoint: "/orderDetails/get",
			method:   methodGet,
		},
		{
			endpoint: "/orderDetails/update",
			method:   methodPatch,
		},
		{
			endpoint: "/orders/orders",
			method:   methodGet,
		},
		{
			endpoint: "/payments/delete",
			method:   methodDelete,
		},
		{
			endpoint: "/payments/get",
			method:   methodGet,
		},
		{
			endpoint: "/payments/payments",
			method:   methodGet,
		},
		{
			endpoint: "/payments/update",
			method:   methodPatch,
		},
		{
			endpoint: "/productLines/productLines",
			method:   methodGet,
		},
		{
			endpoint: "/products/delete",
			method:   methodDelete,
		},
		{
			endpoint: "/products/get",
			method:   methodGet,
		},
		{
			endpoint: "/products/update",
			method:   methodPatch,
		},
	}

	codes := []string{codeOk, codeBadRequest, codeInternalServerError, codeServiceUnavailable}

	for _, l := range apiRequestCountLabels {
		for _, code := range codes {
			for _, originator := range originators {
				RequestCount.WithLabelValues(host, code, l.method, l.endpoint, originator).Add(0)
			}
		}
	}
}

// DB metrics counters

func initializeDbMetricsCounters() {
	dbFunctionLabels := []string{
		"repository.CustomerQueryRepository.GetCustomers",
		"repository.CustomerQueryRepository.FindCustomer",
		"repository.CustomerQueryRepository.DeleteCustomer",
		"repository.CustomerQueryRepository.UpdateCustomer",
		"repository.EmployeeQueryRepository.GetEmployees",
		"repository.EmployeeQueryRepository.FindEmployee",
		"repository.EmployeeQueryRepository.DeleteEmployee",
		"repository.EmployeeQueryRepository.UpdateEmployee",
		"repository.OfficeQueryRepository.GetOffices",
		"repository.OfficeQueryRepository.FindOffice",
		"repository.OfficeQueryRepository.DeleteOffice",
		"repository.OfficeQueryRepository.UpdateOffice",
		"repository.OrderDetailQueryRepository.GetOrderDetails",
		"repository.OrderDetailQueryRepository.FindOrderDetails",
		"repository.OrderDetailQueryRepository.UpdateOrderDetails",
		"repository.OrderDetailQueryRepository.DeleteOrder",
		"repository.OrderQueryRepository.Getorders",
		"repository.PaymentQueryRepository.GetPayments",
		"repository.PaymentQueryRepository.FindPayment",
		"repository.PaymentQueryRepository.DeletePayment",
		"repository.PaymentQueryRepository.UpdatePayment",
		"repository.ProductQueryRepository.GetProducts",
		"repository.ProductQueryRepository.FindProduct",
		"repository.ProductQueryRepository.DeleteProduct",
		"repository.ProductQueryRepository.UpdateProduct",
		"repository.ProductLineQueryRepository.GetProductLines",
	}

	for _, functionLabel := range dbFunctionLabels {
		DbCount.WithLabelValues(host, functionLabel).Add(0)
	}

	DbError.WithLabelValues(host).Add(0)
}
