package dbmod


type FileConf struct {
	FileConfId		int64 		`xorm:"pk autoincr 'file_conf_id'" orm:"pk"`
	Siteid			int64		`xorm:"'siteid'"`
	FileCategory 	int			`xorm:"'file_category'"`
	FileType 		int			`xorm:"'file_type'"`
	StorageType 	int			`xorm:"'storage_type'"`
	Status 			int			`xorm:"'status'"`
	FilePath 		string		`xorm:"'file_path'"`
	FilePrefix 		string		`xorm:"'file_prefix'"`
	Template 		string		`xorm:"'template'"`
	CreateTime 		string		`xorm:"'create_time'"`
	Remark 			string		`xorm:"text 'remark'"`
}

type FileInfo struct {
	Fileid			int64 		`xorm:"pk autoincr 'fileid'" orm:"pk"`
	Siteid			int64		`xorm:"'siteid'"`
	Filekey			string		`xorm:"'filekey'"`
	Name 			string		`xorm:"'name'"`
	OrigUrl			string		`xorm:"varchar(512) notnull 'orig_url'"`
	LocalFile		string		`xorm:"varchar(512) notnull 'local_file'"`
	Redirect		string		`xorm:"varchar(512) 'redirect'"`
	Location		string
	FileOwner		string		`xorm:"varchar(255) notnull 'file_owner'"`
	FileNo 			string		`xorm:"varchar(255) notnull 'file_no'"`
	FileSize		int64		`xorm:"bigint notnull 'file_size'"`
	Hash			string		`xorm:"'hash'"`
	StorageType 	int			`xorm:"'storage_type'"`
	Status 			int			`xorm:"'status'"`
	CreateTime 		string		`xorm:"'create_time'"`
	Remark 			string		`xorm:"text 'remark'"`
}

type FileStorage struct {
	Fileid			int64 		`xorm:"pk 'fileid'"`
	Status 			int
	Content			[]uint8		`xorm:"MEDIUMBLOB"`
}