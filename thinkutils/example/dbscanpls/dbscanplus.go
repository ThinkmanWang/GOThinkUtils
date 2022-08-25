package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"runtime"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

type ConfigUrl struct {
	Id          thinkutils.NullInt64  `json:"id" field:"id"`
	ShowName    thinkutils.NullString `json:"showName" field:"show_name"`
	Name        thinkutils.NullString `json:"name" field:"name"`
	ChannelName thinkutils.NullString `json:"channelName" field:"channel_name"`
	PackageName thinkutils.NullString `json:"packageName" field:"package_name"`
	//ConfigUrl    thinkutils.NullString `json:"configUrl" field:"config_url"`
	//ShowCallback thinkutils.NullString `json:"showCallback" field:"show_callback"`
}

func scanTest1() {
	db := thinkutils.ThinkMysql.QuickConn()
	rows, err := db.Query(`
		SELECT
			b.id
			, b.show_name
			, b.name
			, c.name as channel_name
			, a.package_name
			, CONCAT('https://openapi.shiyculture.com/SDKBase/getGameInitConfig?gameId=',b.name,'&channelId=', c.name) as config_url
			, CONCAT('https://openapi.shiyculture.com/SDKBase/ad-callback/',c.name,'/', b.name, '/show?os=$os$&os_version=$ov$&model=$m$&lang=$lan$&country=$c$&width=$w$&height=$h$&pkg=$pkg$&app_version=$av$&useragent=$ua$&referer=$rf$&net_type=$nt$&carrier=$ca$&progress=$progress$&imei=$im$&oaid=__OAID__&ad_id=$ad$&ad_name=$an$&req_id=$req$') as show_callback
			, CONCAT('https://openapi.shiyculture.com/SDKBase/ad-callback/',c.name,'/', b.name, '/click?os=$os$&os_version=$ov$&model=$m$&lang=$lan$&country=$c$&width=$w$&height=$h$&pkg=$pkg$&app_version=$av$&useragent=$ua$&referer=$rf$&net_type=$nt$&carrier=$ca$&progress=$progress$&imei=$im$&oaid=__OAID__&ad_id=$ad$&ad_name=$an$&req_id=$req$') as click_callback
			, CONCAT('https://openapi.shiyculture.com/SDKBase/ad-callback/',c.name,'/', b.name, '/active?oaid=654321&imei=123456') as active_callback
		FROM
			t_game_package_name as a
			left join t_game AS b on a.game_id = b.id
			left join t_channel as c on a.channel_id = c.id
		ORDER BY
			b.id
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		game := &ConfigUrl{}

		err = thinkutils.ThinkMysql.ScanRow(rows, game)
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info("%s", thinkutils.JSONUtils.ToJson(game))
	}
}

func scructTest() {
	config := ConfigUrl{}
	szName, bExists := thinkutils.StructUtils.FieldNameByTag(&config, "field", "id")
	log.Info("%o => %s", bExists, szName)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Info("Hello World")

	//scructTest()
	scanTest1()
}
