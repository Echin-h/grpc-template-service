package conf

type GlobalConfig struct {
	MODE string `yaml:"Mode"`
	Host string `yaml:"Host"`
	Port string `yaml:"Port"` // grpc和http服务监听端口
	Log  struct {
		LogPath string `yaml:"LogPath"`
		CLS     struct {
			Endpoint    string `yaml:"Endpoint"`
			AccessKey   string `yaml:"AccessKey"`
			AccessToken string `yaml:"AccessToken"`
			TopicID     string `yaml:"TopicID"`
		} `yaml:"CLS"`
	} `yaml:"Log"`
	Postgres struct {
		Addr     string `yaml:"Addr"`
		PORT     string `yaml:"Port"`
		USER     string `yaml:"User"`
		PASSWORD string `yaml:"Password"`
		DATABASE string `yaml:"Database"`
		UseTLS   bool   `yaml:"UseTLS"`
		Debug    bool   `yaml:"Debug"`
	} `yaml:"Postgres"`
	Redis struct {
		Addr     string `yaml:"Addr"`
		PORT     string `yaml:"Port"`
		PASSWORD string `yaml:"Password"`
		DB       int    `yaml:"Db"`
	} `yaml:"Redis"`
	B2 struct {
		BucketKeyID string `yaml:"BucketKeyID"`
		BucketKey   string `yaml:"BucketKey"`
		BucketName  string `yaml:"BucketName"`
	} `yaml:"B2"`
	SentryDsn string `yaml:"SentryDsn"`
	HDUHELP   struct {
		ClientID        string `yaml:"ClientID"`
		ClientSecret    string `yaml:"ClientSecret"`
		OauthFinishPath string `yaml:"OauthFinishPath"`
	}
	Otel struct {
		Enable         bool   `yaml:"Enable"`
		ServiceName    string `yaml:"ServiceName"`
		ServiceVersion string `yaml:"ServiceVersion"`
		AgentHost      string `yaml:"AgentHost"`
		AgentPort      string `yaml:"AgentPort"`
	}
}
