package cyw4343w

const (
	maxCapsBufferSize = 768
	busHeaderLen      = 12
	ioctlOffset       = 4 + 12 + 16 // sizeof(void*) + 12 + 16

	cdcfIocError   = 0x01       // 0=success, 1=ioctl cmd failed.
	cdcfIocIfMask  = 0xF000     // I/F index.
	cdcfIocIfShift = 12         // # of bits of shift for I/F Mask.
	cdcfIocIdMask  = 0xFFFF0000 // used to uniquely id an ioctl req/resp pairing.
	cdcfIocIdShift = 16         // # of bits of shift for ID Mask.
)

const (
	iovarStrBtaddr               = "bus:btsdiobufaddr"
	iovarStrActframe             = "actframe"
	iovarStrBss                  = "bss"
	iovarStrBssRateset           = "bss_rateset"
	iovarStrCsa                  = "csa"
	iovarStrAmpduTid             = "ampdu_tid"
	iovarStrApsta                = "apsta"
	iovarStrAllmulti             = "allmulti"
	iovarStrCountry              = "country"
	iovarStrEventMsgs            = "event_msgs"
	iovarStrEventMsgsExt         = "event_msgs_ext"
	iovarStrEscan                = "escan"
	iovarStrSupWpa               = "sup_wpa"
	iovarStrCurEtheraddr         = "cur_etheraddr"
	iovarStrQtxpower             = "qtxpower"
	iovarStrMcastList            = "mcast_list"
	iovarStrPm2SleepRet          = "pm2_sleep_ret"
	iovarStrPmLimit              = "pm_limit"
	iovarStrListenIntervalBeacon = "bcn_li_bcn"
	iovarStrListenIntervalDtim   = "bcn_li_dtim"
	iovarStrListenIntervalAssoc  = "assoc_listen"
	iovarPspollPeriod            = "pspoll_prd"
	iovarStrVendorIe             = "vndr_ie"
	iovarStrTxGlom               = "bus:txglom"
	iovarStrActionFrame          = "actframe"
	iovarStrAcParamsSta          = "wme_ac_sta"
	iovarStrCounters             = "counters"
	iovarStrPktFilterAdd         = "pkt_filter_add"
	iovarStrPktFilterDelete      = "pkt_filter_delete"
	iovarStrPktFilterEnable      = "pkt_filter_enable"
	iovarStrPktFilterMode        = "pkt_filter_mode"
	iovarStrPktFilterList        = "pkt_filter_list"
	iovarStrPktFilterStats       = "pkt_filter_stats"
	iovarStrPktFilterClearStats  = "pkt_filter_clear_stats"
	iovarStrDutyCycleCck         = "dutycycle_cck"
	iovarStrDutyCycleOfdm        = "dutycycle_ofdm"
	iovarStrMkeepAlive           = "mkeep_alive"
	iovarStrVersion              = "ver"
	iovarStrSupWpa2Eapver        = "sup_wpa2_eapver"
	iovarStrRoamOff              = "roam_off"
	iovarStrClosednet            = "closednet"
	iovarStrP2pDisc              = "p2p_disc"
	iovarStrP2pDev               = "p2p_dev"
	iovarStrP2pIf                = "p2p_if"
	iovarStrP2pIfadd             = "p2p_ifadd"
	iovarStrP2pIfdel             = "p2p_ifdel"
	iovarStrP2pIfupd             = "p2p_ifupd"
	iovarStrP2pScan              = "p2p_scan"
	iovarStrP2pState             = "p2p_state"
	iovarStrP2pSsid              = "p2p_ssid"
	iovarStrP2pIpAddr            = "p2p_ip_addr"
	iovarStrNrate                = "nrate"
	iovarStrBgrate               = "bg_rate"
	iovarStrArate                = "a_rate"
	iovarStrNmode                = "nmode"
	iovarStrMaxAssoc             = "maxassoc"
	iovarStr2gMulticastRate      = "2g_mrate"
	iovarStr2gRate               = "2g_rate"
	iovarStrMpc                  = "mpc"
	iovarStrIbssJoin             = "IBSS_join_only"
	iovarStrAmpduBaWindowSize    = "ampdu_ba_wsize"
	iovarStrAmpduMpdu            = "ampdu_mpdu"
	iovarStrAmpduRx              = "ampdu_rx"
	iovarStrAmpduRxFactor        = "ampdu_rx_factor"
	iovarStrAmpduHostReorder     = "ampdu_hostreorder"
	iovarStrMimoBwCap            = "mimo_bw_cap"
	iovarStrRmcAckreq            = "rmc_ackreq"
	iovarStrRmcStatus            = "rmc_status"
	iovarStrRmcCounts            = "rmc_stats"
	iovarStrRmcRole              = "rmc_role"
	iovarStrHt40Intolerance      = "intol40"
	iovarStrRand                 = "rand"
	iovarStrSsid                 = "ssid"
	iovarStrWsec                 = "wsec"
	iovarStrWpaAuth              = "wpa_auth"
	iovarStrInterfaceRemove      = "interface_remove"
	iovarStrSupWpaTmo            = "sup_wpa_tmo"
	iovarStrJoin                 = "join"
	iovarStrTlv                  = "tlv"
	iovarStrNphyAntsel           = "nphy_antsel"
	iovarStrAvbTimestampAddr     = "avb_timestamp_addr"
	iovarStrBssMaxAssoc          = "bss_maxassoc"
	iovarStrRmReq                = "rm_req"
	iovarStrRmRep                = "rm_rep"
	iovarStrPspretendRetryLimit  = "pspretend_retry_limit"
	iovarStrPspretendThreshold   = "pspretend_threshold"
	iovarStrSwdivTimeout         = "swdiv_timeout"
	iovarStrResetCnts            = "reset_cnts"
	iovarStrPhyrateLog           = "phyrate_log"
	iovarStrPhyrateLogSize       = "phyrate_log_size"
	iovarStrPhyrateLogDump       = "phyrate_dump"
	iovarStrScanAssocTime        = "scan_assoc_time"
	iovarStrScanUnassocTime      = "scan_unassoc_time"
	iovarStrScanPassiveTime      = "scan_passive_time"
	iovarStrScanHomeTime         = "scan_home_time"
	iovarStrScanNprobes          = "scan_nprobes"
	iovarStrAutocountry          = "autocountry"
	iovarStrCap                  = "cap"
	iovarStrMpduPerAmpdu         = "ampdu_mpdu"
	iovarStrVhtFeatures          = "vht_features"
	iovarStrChanspec             = "chanspec"
	iovarStrMgmtFrame            = "mgmt_frame"
	iovarStrWowl                 = "wowl"
	iovarStrWowlOs               = "wowl_os"
	iovarStrWowlActivate         = "wowl_activate"
	iovarStrWowlClear            = "wowl_clear"
	iovarStrWowlActivateSecure   = "wowl_activate_secure"
	iovarStrWowlSecSessInfo      = "wowl_secure_sess_info"
	iovarStrWowlKeepAlive        = "wowl_keepalive"
	iovarStrWowlPattern          = "wowl_pattern"
	iovarStrWowlPatternClr       = "clr"
	iovarStrWowlPatternAdd       = "add"
	iovarStrWowlArpHostIp        = "wowl_arp_hostip"
	iovarStrUlpWait              = "ulp_wait"
	iovarStrUlp                  = "ulp"
	iovarStrUlpHostIntrMode      = "ulp_host_intr_mode"
	iovarStrDump                 = "dump"
	iovarStrPnoOn                = "pfn"
	iovarStrPnoAdd               = "pfn_add"
	iovarStrPnoSet               = "pfn_set"
	iovarStrPnoClear             = "pfnclear"
	iovarStrScanCacheClear       = "scancache_clear"
	mcsSetlen                    = 16
	iovarStrRrm                  = "rrm"
	iovarStrRrmNoiseReq          = "rrm_noise_req"
	iovarStrRrmNbrReq            = "rrm_nbr_req"
	iovarStrRrmLmReq             = "rrm_lm_req"
	iovarStrRrmStatReq           = "rrm_stat_req"
	iovarStrRrmFrameReq          = "rrm_frame_req"
	iovarStrRrmChloadReq         = "rrm_chload_req"
	iovarStrRrmBcnReq            = "rrm_bcn_req"
	iovarStrRrmNbrList           = "rrm_nbr_list"
	iovarStrRrmNbrAdd            = "rrm_nbr_add_nbr"
	iovarStrRrmNbrDel            = "rrm_nbr_del_nbr"
	iovarStrRrmBcnreqThrtlWin    = "rrm_bcn_req_thrtl_win"
	iovarStrRrmBcnreqMaxoffTime  = "rrm_bcn_req_max_off_chan_time"
	iovarStrRrmBcnreqTrfmsPrd    = "rrm_bcn_req_traff_meas_per"
	iovarStrWnm                  = "wnm"
	iovarStrBsstransQuery        = "wnm_bsstrans_query"
	iovarStrBsstransResp         = "wnm_bsstrans_resp"
	iovarStrMeshAddRoute         = "mesh_add_route"
	iovarStrMeshDelRoute         = "mesh_del_route"
	iovarStrMeshFind             = "mesh_find"
	iovarStrMeshFilter           = "mesh_filter"
	iovarStrMeshPeer             = "mesh_peer"
	iovarStrMeshPeerStatus       = "mesh_peer_status"
	iovarStrMeshDelfilter        = "mesh_delfilter"
	iovarStrMeshMaxPeers         = "mesh_max_peers"
	iovarStrFbtOverDs            = "fbtoverds"
	iovarStrFbtCapabilities      = "fbt_cap"
	iovarStrMfp                  = "mfp"
	iovarStrBip                  = "bip"
	iovarStrOtpraw               = "otpraw"
	iovarNan                     = "nan"
	iovarStrClmload              = "clmload"
	iovarStrClmloadStatus        = "clmload_status"
	iovarStrClmver               = "clmver"
	iovarStrMemuse               = "memuse"
	iovarStrLdpcCap              = "ldpc_cap"
	iovarStrLdpcTx               = "ldpc_tx"
	iovarStrSgiRx                = "sgi_rx"
	iovarStrSgiTx                = "sgi_tx"
	iovarStrApivtwOverride       = "brcmapivtwo"
	iovarStrBwteBwteGciMask      = "bwte_gci_mask"
	iovarStrBwteGciSendmsg       = "bwte_gci_sendm"
	iovarStrWdDisable            = "wd_disable"
	iovarStrDltro                = "dltro"
	iovarStrSaePassword          = "sae_password"
	iovarStrSaePweLoop           = "sae_max_pwe_loop"
	iovarStrPmkidInfo            = "pmkid_info"
	iovarStrPmkidClear           = "pmkid_clear"
	iovarStrAuthStatus           = "auth_status"
	iovarStrBtcLescanParams      = "btc_lescan_params"
	iovarStrArpVersion           = "arp_version"
	iovarStrArpPeerage           = "arp_peerage"
	iovarStrArpoe                = "arpoe"
	iovarStrArpOl                = "arp_ol"
	iovarStrArpTableClear        = "arp_table_clear"
	iovarStrArpHostip            = "arp_hostip"
	iovarStrArpHostipClear       = "arp_hostip_clear"
	iovarStrArpStats             = "arp_stats"
	iovarStrArpStatsClear        = "arp_stats_clear"
	iovarStrTko                  = "tko"
	iovarStrRoamTimeThresh       = "roam_time_thresh"
	iovarWnmMaxidle              = "wnm_maxidle"
	iovarStrHe                   = "he"
	iovarStrTwt                  = "twt"
	iovarStrOffloadConfig        = "offload_config"
	iovarStrWsecInfo             = "wsec_info"
	iovarStrKeepaliveConfig      = "keep_alive"
	iovarStrMbo                  = "mbo"
)

const (
	wlcGetMagic = iota
	wlcGetVersion
	wlcUp
	wlcDown
	wlcGetLoop
	wlcSetLoop
	wlcDump
	wlcGetMsglevel
	wlcSetMsglevel
	wlcGetPromisc
	wlcSetPromisc
	wlcGetRate
	wlcGetInstance
	wlcGetInfra
	wlcSetInfra
	wlcGetAuth
	wlcSetAuth
	wlcGetBssid
	wlcSetBssid
	wlcGetSsid
	wlcSetSsid
	wlcRestart
	wlcGetChannel
	wlcSetChannel
	wlcGetSrl
	wlcSetSrl
	wlcGetLrl
	wlcSetLrl
	wlcGetPlcphdr
	wlcSetPlcphdr
	wlcGetRadio
	wlcSetRadio
	wlcGetPhytype
	wlcDumpRate
	wlcSetRateParams
	wlcGetKey
	wlcSetKey
	wlcGetRegulatory
	wlcSetRegulatory
	wlcGetPassiveScan
	wlcSetPassiveScan
	wlcScan
	wlcScanResults
	wlcDisassoc
	wlcReassoc
	wlcGetRoamTrigger
	wlcSetRoamTrigger
	wlcGetRoamDelta
	wlcSetRoamDelta
	wlcGetRoamScanPeriod
	wlcSetRoamScanPeriod
	wlcEvm
	wlcGetTxant
	wlcSetTxant
	wlcGetAntdiv
	wlcSetAntdiv
	wlcGetClosed
	wlcSetClosed
	wlcGetMaclist
	wlcSetMaclist
	wlcGetRateset
	wlcSetRateset
	wlcLongtrain
	wlcGetBcnprd
	wlcSetBcnprd
	wlcGetDtimprd
	wlcSetDtimprd
	wlcGetSrom
	wlcSetSrom
	wlcGetWepRestrict
	wlcSetWepRestrict
	wlcGetCountry
	wlcSetCountry
	wlcGetPm
	wlcSetPm
	wlcGetWake
	wlcSetWake
	wlcGetForcelink
	wlcSetForcelink
	wlcFreqAccuracy
	wlcCarrierSuppress
	wlcGetPhyreg
	wlcSetPhyreg
	wlcGetRadioreg
	wlcSetRadioreg
	wlcGetRevinfo
	wlcGetUcantdiv
	wlcSetUcantdiv
	wlcRReg
	wlcWReg
	wlcGetMacmode
	wlcSetMacmode
	wlcGetMonitor
	wlcSetMonitor
	wlcGetGmode
	wlcSetGmode
	wlcGetLegacyErp
	wlcSetLegacyErp
	wlcGetRxAnt
	wlcGetCurrRateset
	wlcGetScansuppress
	wlcSetScansuppress
	wlcGetAp
	wlcSetAp
	wlcGetEapRestrict
	wlcSetEapRestrict
	wlcScbAuthorize
	wlcScbDeauthorize
	wlcGetWdslist
	wlcSetWdslist
	wlcGetAtim
	wlcSetAtim
	wlcGetRssi
	wlcGetPhyantdiv
	wlcSetPhyantdiv
	wlcApRxOnly
	wlcGetTxPathPwr
	wlcSetTxPathPwr
	wlcGetWsec
	wlcSetWsec
	wlcGetPhyNoise
	wlcGetBssInfo
	wlcGetPktcnts
	wlcGetLazywds
	wlcSetLazywds
	wlcGetBandlist
	wlcGetBand
	wlcSetBand
	wlcScbDeauthenticate
	wlcGetShortslot
	wlcGetShortslotOverride
	wlcSetShortslotOverride
	wlcGetShortslotRestrict
	wlcSetShortslotRestrict
	wlcGetGmodeProtection
	wlcGetGmodeProtectionOverride
	wlcSetGmodeProtectionOverride
	wlcUpgrade
	wlcGetIgnoreBcns
	wlcSetIgnoreBcns
	wlcGetScbTimeout
	wlcSetScbTimeout
	wlcGetAssoclist
	wlcGetClk
	wlcSetClk
	wlcGetUp
	wlcOut
	wlcGetWpaAuth
	wlcSetWpaAuth
	wlcGetUcflags
	wlcSetUcflags
	wlcGetPwridx
	wlcSetPwridx
	wlcGetTssi
	wlcGetSupRatesetOverride
	wlcSetSupRatesetOverride
	wlcGetProtectionControl
	wlcSetProtectionControl
	wlcGetPhylist
	wlcEncryptStrength
	wlcDecryptStatus
	wlcGetKeySeq
	wlcGetScanChannelTime
	wlcSetScanChannelTime
	wlcGetScanUnassocTime
	wlcSetScanUnassocTime
	wlcGetScanHomeTime
	wlcSetScanHomeTime
	wlcGetScanNprobes
	wlcSetScanNprobes
	wlcGetPrbRespTimeout
	wlcSetPrbRespTimeout
	wlcGetAtten
	wlcSetAtten
	wlcGetShmem
	wlcSetShmem
	wlcSetWsecTest
	wlcScbDeauthenticateForReason
	wlcTkipCountermeasures
	wlcGetPiomode
	wlcSetPiomode
	wlcSetAssocPrefer
	wlcGetAssocPrefer
	wlcSetRoamPrefer
	wlcGetRoamPrefer
	wlcSetLed
	wlcGetLed
	wlcGetInterferenceMode
	wlcSetInterferenceMode
	wlcGetChannelQa
	wlcStartChannelQa
	wlcGetChannelSel
	wlcStartChannelSel
	wlcGetValidChannels
	wlcGetFakefrag
	wlcSetFakefrag
	wlcGetPwroutPercentage
	wlcSetPwroutPercentage
	wlcSetBadFramePreempt
	wlcGetBadFramePreempt
	wlcSetLeapList
	wlcGetLeapList
	wlcGetCwmin
	wlcSetCwmin
	wlcGetCwmax
	wlcSetCwmax
	wlcGetWet
	wlcSetWet
	wlcGetPub
	wlcGetKeyPrimary
	wlcSetKeyPrimary
	wlcGetAciArgs
	wlcSetAciArgs
	wlcUnsetCallback
	wlcSetCallback
	wlcGetRadar
	wlcSetRadar
	wlcSetSpectManagment
	wlcGetSpectManagment
	wlcWdsGetRemoteHwaddr
	wlcWdsGetWpaSup
	wlcSetCsScanTimer
	wlcGetCsScanTimer
	wlcMeasureRequest
	wlcInit
	wlcSendQuiet
	wlcKeepalive
	wlcSendPwrConstraint
	wlcUpgradeStatus
	wlcCurrentPwr
	wlcGetScanPassiveTime
	wlcSetScanPassiveTime
	wlcLegacyLinkBehavior
	wlcGetChannelsInCountry
	wlcGetCountryList
	wlcGetVar
	wlcSetVar
	wlcNvramGet
	wlcNvramSet
	wlcNvramDump
	wlcReboot
	wlcSetWsecPmk
	wlcGetAuthMode
	wlcSetAuthMode
	wlcGetWakeentry
	wlcSetWakeentry
	wlcNdconfigItem
	wlcNvotpw
	wlcOtpw
	wlcIovBlockGet
	wlcIovModulesGet
	wlcSoftReset
	wlcGetAllowMode
	wlcSetAllowMode
	wlcGetDesiredBssid
	wlcSetDesiredBssid
	wlcDisassocMyap
	wlcGetNbands
	wlcGetBandstates
	wlcGetWlcBssInfo
	wlcGetAssocInfo
	wlcGetOidPhy
	wlcSetOidPhy
	wlcSetAssocTime
	wlcGetDesiredSsid
	wlcGetChanspec
	wlcGetAssocState
	wlcSetPhyState
	wlcGetScanPending
	wlcGetScanreqPending
	wlcGetPrevRoamReason
	wlcSetPrevRoamReason
	wlcGetBandstatesPi
	wlcGetPhyState
	wlcGetBssWpaRsn
	wlcGetBssWpa2Rsn
	wlcGetBssBcnTs
	wlcGetIntDisassoc
	wlcSetNumPeers
	wlcGetNumBss
	wlcGetWsecPmk
	wlcGetRandomBytes
	wlcLast
)
