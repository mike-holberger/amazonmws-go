# amazonmws-go
API Client library for Amazon MWS including Orders, Reports, Feeds, etc. APIs

Initialize Client with Credentials:

    creds := amazonmwsapi.Creds{}
    mustMapEnv(&creds.AccessID, "ACCESS_ID", "")
    mustMapEnv(&creds.AccessKey, "ACCESS_KEY", "")
    mustMapEnv(&creds.Merchant, "MERCHANT_ID", "")
    amazonClient = amazonmwsapi.NewAmazonClient(creds, "US", nil)

    func mustMapEnv(target *string, envKey string, useDefault string) {
	    v := os.Getenv(envKey)
	    if v == "" {
	        v = useDefault
	    }
	        *target = v
    }
    
    
Initialize APIs

    ordersAPI := amazonmwsapi.NewOrdersAPI(amazonClient)
    feedsAPI := amazonmwsapi.NewFeedsAPI(amazonClient)
    reportsAPI := amazonmwsapi.NewReportsAPI(amazonClient)
