package orm

type (
	// target
	DMTargetResource struct {
		Id int64
		Description string
		MD5 string `xorm:"unique"`
		Path string `xorm:"unique"`
		TagIdList []int64
		BackupIdList []int64
		ChildId int64
		ChildGenre string
	}
	DMBinaryResource struct {
		Id int64
		ParentId int64
		CompatibilityDescription string
	}
	DMImageResource struct {
		Id int64
		ParentId int64
	}
	DMMusicResource struct {
		Id int64
		ParentId int64
		Commentary string
	}
	DMVideoResource struct {
		Id int64
		ParentId int64
	}
	// backup
	DMBackupResource struct {
		Id int64
		Path string `xorm:"unique"`
		// 例如一个备份一个文件夹，则其中的所有文件都会被该备份资源所备份
		BackupTargetIds[] int64
	}
	// tag
	DMTag struct {
		Id int64
		Name string
	}
	// permission
	DMWhiteList struct {
		Id int64
		Level int `xorm:"unique"`
		AccountIds [] int64
	}
)

func ( PR *DMTargetResource ) GetChildResource() ( chR *interface{}, er error ) {
	if has, e := engine_dm.Table("d_m_" + PR.ChildGenre + "_resource" ).ID(PR.ChildId).Get(chR); !has || e != nil {
		er = e
		return
	}
	return
}

