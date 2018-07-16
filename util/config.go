package util

//一个view的名称，对应一组cidr设置
//匹配的时候，根据输入的一个IP，从配置中获取对应的view_id，如果没有匹配成功，则id为-1
//同理，根据输入的一个域名name，可以从配置中获取对应的name_id，
//最终的fwd是根据name_id，view_id联合确定的，因为fwd的数量不会很大，我们可以在这里用简单遍历的方式对fwd进行匹配
//（如果很大，那么数据结构上应该再进行调整）使用map[[2]int],int的方式来进行快速匹配

type cdn_ips struct {
	name2id	map[string]int
	cidr []dns_cidr
}