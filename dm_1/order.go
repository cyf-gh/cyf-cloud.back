// 进行归档时候的操作
package dm_1

// 递归某个目录的所有文件
// 结果不包含 r 自身
// 不会递归二进制文件夹，见 IsBinaryDirectory1()
func DMRecruitLs( r DMResource ) ( rs []DMResource ) {
	rs = []DMResource{}

	if r.IsDire() {
		if r.IsBinaryDirectory1() {
			return
		}
		rrs, _ := r.Ls()
		for _, rr := range rrs {
			rs = append(rs, rr)
			if rsr := DMRecruitLs( rr ); len(rsr) != 0 {
				rs = append(rs, rsr...)
			}
		}
	}
	return
}
