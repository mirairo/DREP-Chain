package blockmgr

import (
	dlog "github.com/drep-project/drep-chain/pkgs/log"
)

const (
	maxHeaderHashCountReq = 255 //最多一次请求的头部hash个数
	maxBlockCountReq      = 16  //最多一次请求的头部hash个数
	maxSyncSleepTime      = 200 //同步的过程中，每个周期休息200毫秒
	maxNetworkTimeout     = 10  //最大网络超时时间
	maxLivePeer           = 20
	broadcastRatio        = 3    //非本地产生的消息，广播的个数是broadcastRatio分之一
	maxTxsCount           = 1024 //最多一次传输交易的个数
	pendingTimerCount     = 10   //同步区块时，最多同时并发的获取块请求的协程数目

	MODULENAME = "blockmgr"
)

var (
	log = dlog.NewLogger(MODULENAME)
)
