package mipush

const (
	HostDevbase = "https://sandbox.xmpush.xiaomi.com"
	HostProBase = "https://api.xmpush.xiaomi.com"
	// auth prefix
	AUTH_PREFIX = "key="
	// “result”: string，”ok” 表示成功, “error” 表示失败。
	RESULT_OK    = "ok"
	RESULT_ERROR = "error"
	comma        = ","
	//	topic列表，使用;$;分割。注: topics参数需要和topic_op参数配合使用，另外topic的数量不能超过5。参数仅适用于“/message/multi_topic”HTTP API。
	multi_topic_split = ";$;"

	Reg_url                          = "/v3/message/regid"
	alias_url                        = "/v3/message/alias"
	User_Account_Url                 = "/v2/message/user_account"
	TOPIC_URL                        = "/v3/message/topic"
	MULTI_TOPIC_URL                  = "/v3/message/multi_topic"
	all_url                          = "/v3/message/all"
	multi_messages_regids_url        = "/v2/multi_messages/regids"
	multi_messages_aliases_url       = "/v2/multi_messages/aliases"
	Multi_Messages_User_Accounts_Url = "/v2/multi_messages/user_accounts"
	stats_url                        = "/v1/stats/message/counters"
	message_trace_url                = "/v1/trace/message/status"
	messages_trace_url               = "/v1/trace/messages/status"
	validation_regids_url            = "/v1/validation/regids"
	subscribe_url                    = "/v2/topic/subscribe"
	unsubscribe_url                  = "/v2/topic/unsubscribe"
	subscribe_alias_url              = "/v2/topic/subscribe/alias"
	unsubscribe_alias_url            = "/v2/topic/unsubscribe/alias"
	fetch_invalid_regids_url         = "https://feedback.xmpush.xiaomi.com/v1/feedback/fetch_invalid_regids"
	delete_schedule_job              = "/v2/schedule_job/delete"
	check_schedule_job_exist         = "/v2/schedule_job/exist"
	get_all_aliases                  = "/v1/alias/all"
	get_all_topics                   = "/v1/topic/all"

	// topic之间的操作关系。支持以下三种：TOPIC_OP_
	TOPIC_OP_UNION        = "UNION"        // 并集
	TOPIC_OP_INTERSECTION = "INTERSECTION" // 交集
	TOPIC_OP_EXCEPT       = "EXCEPT"       // 差集

	// notify_type的值可以是DEFAULT_ALL或者以下其他几种的OR组合：NOTIFY_TYPE_
	NOTIFY_TYPE_DEFAULT_ALL     = -1
	NOTIFY_TYPE_DEFAULT_SOUND   = 1 //  使用默认提示音提示；
	NOTIFY_TYPE_DEFAULT_VIBRATE = 2 //  使用默认震动提示
	NOTIFY_TYPE_DEFAULT_LIGHTS  = 4 //  使用默认led灯光提示；
)
