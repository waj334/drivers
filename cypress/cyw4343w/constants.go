package cyw4343w

const (
	busHeaderLen = 12
	ioctlOffset  = 4 + 12 + 16 // sizeof(void*) + 12 + 16
	ioctlMaxlen

	dataHeader       = 2
	asynceventHeader = 1
	controlHeader    = 0

	// BDC (Broadcom Dongle Communication) header flags.
	bdcVersion      = 2
	bdcVersionShift = 4
	bdcFlagVersion  = bdcVersion << bdcVersionShift // 0x20

	cdcfIocError   = 0x01       // 0=success, 1=ioctl cmd failed.
	cdcfIocIfMask  = 0xF000     // I/F index.
	cdcfIocIfShift = 12         // # of bits of shift for I/F Mask.
	cdcfIocIdMask  = 0xFFFF0000 // used to uniquely id an ioctl req/resp pairing.
	cdcfIocIdShift = 16         // # of bits of shift for ID Mask.

	maxBackplaneTransferSize = 64 * 24
	backplaneReadPaddSize    = 0
	backplaneAddressMask     = 0x7FFF
	backplaneWindowSize      = backplaneAddressMask + 1

	sbOftAddrMask   = 0x07FFF
	sbOftAddrLimit  = 0x08000
	sbAccess24BFlag = 0x08000

	platformWlanRamBase = 0x0
	wlanBusUpAttempts   = 1000
	htAvailWaitMs       = 1
	ksoWaitMs           = 1
	ksoWakeMs           = 3
	maxKsoAttempts      = 64
	maxCapsBufferSize   = 768

	aiIoctrlOffset        = 0x408
	sicfFgc               = 0x0002
	sicfClockEn           = 0x0001
	sicfCpuhalt           = 0x0020
	aiResetctrlOffset     = 0x800
	aiResetstatusOffset   = 0x804
	aircReset             = 1
	wrapperRegisterOffset = 0x100000

	nvmImageSizeAlignment = 4
)

const (
	/* Backplane architecture */
	chipcommonBaseAddress = 0x18000000 /* Chipcommon core register region   */
	i2s0BaseAddress       = 0x18001000 /* I2S0 core register region     */
	i2s1BaseAddress       = 0x18002000 /* I2S1 core register region     */
	appsArmcr4BaseAddress = 0x18003000 /* Apps Cortex-R4 core register region     */
	DmaBaseAddress        = 0x18004000 /* DMA core register region     */
	gmacBaseAddress       = 0x18005000 /* GMAC core register region     */
	usb20h0BaseAddress    = 0x18006000 /* USB20H0 core register region     */
	usb20dBaseAddress     = 0x18007000 /* USB20D core register region     */
	sdiohBaseAddress      = 0x18008000 /* SDIOH Device core register region */
	dot11macBaseAddress   = 0x18001000
)

const (
	IovarStrBtaddr               = "bus:btsdiobufaddr"
	IovarStrActframe             = "actframe"
	IovarStrBss                  = "bss"
	IovarStrBssRateset           = "bss_rateset"
	IovarStrCsa                  = "csa"
	IovarStrAmpduTid             = "ampdu_tid"
	IovarStrApsta                = "apsta"
	IovarStrAllmulti             = "allmulti"
	IovarStrCountry              = "country"
	IovarStrEventMsgs            = "event_msgs"
	IovarStrEventMsgsExt         = "event_msgs_ext"
	IovarStrEscan                = "escan"
	IovarStrSupWpa               = "sup_wpa"
	IovarStrCurEtheraddr         = "cur_etheraddr"
	IovarStrQtxpower             = "qtxpower"
	IovarStrMcastList            = "mcast_list"
	IovarStrPm2SleepRet          = "pm2_sleep_ret"
	IovarStrPmLimit              = "pm_limit"
	IovarStrListenIntervalBeacon = "bcn_li_bcn"
	IovarStrListenIntervalDtim   = "bcn_li_dtim"
	IovarStrListenIntervalAssoc  = "assoc_listen"
	iovarPspollPeriod            = "pspoll_prd"
	IovarStrVendorIe             = "vndr_ie"
	IovarStrTxGlom               = "bus:txglom"
	IovarStrActionFrame          = "actframe"
	IovarStrAcParamsSta          = "wme_ac_sta"
	IovarStrCounters             = "counters"
	IovarStrPktFilterAdd         = "pkt_filter_add"
	IovarStrPktFilterDelete      = "pkt_filter_delete"
	IovarStrPktFilterEnable      = "pkt_filter_enable"
	IovarStrPktFilterMode        = "pkt_filter_mode"
	IovarStrPktFilterList        = "pkt_filter_list"
	IovarStrPktFilterStats       = "pkt_filter_stats"
	IovarStrPktFilterClearStats  = "pkt_filter_clear_stats"
	IovarStrDutyCycleCck         = "dutycycle_cck"
	IovarStrDutyCycleOfdm        = "dutycycle_ofdm"
	IovarStrMkeepAlive           = "mkeep_alive"
	IovarStrVersion              = "ver"
	IovarStrSupWpa2Eapver        = "sup_wpa2_eapver"
	IovarStrRoamOff              = "roam_off"
	IovarStrClosednet            = "closednet"
	IovarStrP2pDisc              = "p2p_disc"
	IovarStrP2pDev               = "p2p_dev"
	IovarStrP2pIf                = "p2p_if"
	IovarStrP2pIfadd             = "p2p_ifadd"
	IovarStrP2pIfdel             = "p2p_ifdel"
	IovarStrP2pIfupd             = "p2p_ifupd"
	IovarStrP2pScan              = "p2p_scan"
	IovarStrP2pState             = "p2p_state"
	IovarStrP2pSsid              = "p2p_ssid"
	IovarStrP2pIpAddr            = "p2p_ip_addr"
	IovarStrNrate                = "nrate"
	IovarStrBgrate               = "bg_rate"
	IovarStrArate                = "a_rate"
	IovarStrNmode                = "nmode"
	IovarStrMaxAssoc             = "maxassoc"
	IovarStr2gMulticastRate      = "2g_mrate"
	IovarStr2gRate               = "2g_rate"
	IovarStrMpc                  = "mpc"
	IovarStrIbssJoin             = "IBSS_join_only"
	IovarStrAmpduBaWindowSize    = "ampdu_ba_wsize"
	IovarStrAmpduMpdu            = "ampdu_mpdu"
	IovarStrAmpduRx              = "ampdu_rx"
	IovarStrAmpduRxFactor        = "ampdu_rx_factor"
	IovarStrAmpduHostReorder     = "ampdu_hostreorder"
	IovarStrMimoBwCap            = "mimo_bw_cap"
	IovarStrRmcAckreq            = "rmc_ackreq"
	IovarStrRmcStatus            = "rmc_status"
	IovarStrRmcCounts            = "rmc_stats"
	IovarStrRmcRole              = "rmc_role"
	IovarStrHt40Intolerance      = "intol40"
	IovarStrRand                 = "rand"
	IovarStrSsid                 = "ssid"
	IovarStrWsec                 = "wsec"
	IovarStrWpaAuth              = "wpa_auth"
	IovarStrInterfaceRemove      = "interface_remove"
	IovarStrSupWpaTmo            = "sup_wpa_tmo"
	IovarStrJoin                 = "join"
	IovarStrTlv                  = "tlv"
	IovarStrNphyAntsel           = "nphy_antsel"
	IovarStrAvbTimestampAddr     = "avb_timestamp_addr"
	IovarStrBssMaxAssoc          = "bss_maxassoc"
	IovarStrRmReq                = "rm_req"
	IovarStrRmRep                = "rm_rep"
	IovarStrPspretendRetryLimit  = "pspretend_retry_limit"
	IovarStrPspretendThreshold   = "pspretend_threshold"
	IovarStrSwdivTimeout         = "swdiv_timeout"
	IovarStrResetCnts            = "reset_cnts"
	IovarStrPhyrateLog           = "phyrate_log"
	IovarStrPhyrateLogSize       = "phyrate_log_size"
	IovarStrPhyrateLogDump       = "phyrate_dump"
	IovarStrScanAssocTime        = "scan_assoc_time"
	IovarStrScanUnassocTime      = "scan_unassoc_time"
	IovarStrScanPassiveTime      = "scan_passive_time"
	IovarStrScanHomeTime         = "scan_home_time"
	IovarStrScanNprobes          = "scan_nprobes"
	IovarStrAutocountry          = "autocountry"
	IovarStrCap                  = "cap"
	IovarStrMpduPerAmpdu         = "ampdu_mpdu"
	IovarStrVhtFeatures          = "vht_features"
	IovarStrChanspec             = "chanspec"
	IovarStrMgmtFrame            = "mgmt_frame"
	IovarStrWowl                 = "wowl"
	IovarStrWowlOs               = "wowl_os"
	IovarStrWowlActivate         = "wowl_activate"
	IovarStrWowlClear            = "wowl_clear"
	IovarStrWowlActivateSecure   = "wowl_activate_secure"
	IovarStrWowlSecSessInfo      = "wowl_secure_sess_info"
	IovarStrWowlKeepAlive        = "wowl_keepalive"
	IovarStrWowlPattern          = "wowl_pattern"
	IovarStrWowlPatternClr       = "clr"
	IovarStrWowlPatternAdd       = "add"
	IovarStrWowlArpHostIp        = "wowl_arp_hostip"
	IovarStrUlpWait              = "ulp_wait"
	IovarStrUlp                  = "ulp"
	IovarStrUlpHostIntrMode      = "ulp_host_intr_mode"
	IovarStrDump                 = "dump"
	IovarStrPnoOn                = "pfn"
	IovarStrPnoAdd               = "pfn_add"
	IovarStrPnoSet               = "pfn_set"
	IovarStrPnoClear             = "pfnclear"
	IovarStrScanCacheClear       = "scancache_clear"
	mcsSetlen                    = 16
	IovarStrRrm                  = "rrm"
	IovarStrRrmNoiseReq          = "rrm_noise_req"
	IovarStrRrmNbrReq            = "rrm_nbr_req"
	IovarStrRrmLmReq             = "rrm_lm_req"
	IovarStrRrmStatReq           = "rrm_stat_req"
	IovarStrRrmFrameReq          = "rrm_frame_req"
	IovarStrRrmChloadReq         = "rrm_chload_req"
	IovarStrRrmBcnReq            = "rrm_bcn_req"
	IovarStrRrmNbrList           = "rrm_nbr_list"
	IovarStrRrmNbrAdd            = "rrm_nbr_add_nbr"
	IovarStrRrmNbrDel            = "rrm_nbr_del_nbr"
	IovarStrRrmBcnreqThrtlWin    = "rrm_bcn_req_thrtl_win"
	IovarStrRrmBcnreqMaxoffTime  = "rrm_bcn_req_max_off_chan_time"
	IovarStrRrmBcnreqTrfmsPrd    = "rrm_bcn_req_traff_meas_per"
	IovarStrWnm                  = "wnm"
	IovarStrBsstransQuery        = "wnm_bsstrans_query"
	IovarStrBsstransResp         = "wnm_bsstrans_resp"
	IovarStrMeshAddRoute         = "mesh_add_route"
	IovarStrMeshDelRoute         = "mesh_del_route"
	IovarStrMeshFind             = "mesh_find"
	IovarStrMeshFilter           = "mesh_filter"
	IovarStrMeshPeer             = "mesh_peer"
	IovarStrMeshPeerStatus       = "mesh_peer_status"
	IovarStrMeshDelfilter        = "mesh_delfilter"
	IovarStrMeshMaxPeers         = "mesh_max_peers"
	IovarStrFbtOverDs            = "fbtoverds"
	IovarStrFbtCapabilities      = "fbt_cap"
	IovarStrMfp                  = "mfp"
	IovarStrBip                  = "bip"
	IovarStrOtpraw               = "otpraw"
	iovarNan                     = "nan"
	IovarStrClmload              = "clmload"
	IovarStrClmloadStatus        = "clmload_status"
	IovarStrClmver               = "clmver"
	IovarStrMemuse               = "memuse"
	IovarStrLdpcCap              = "ldpc_cap"
	IovarStrLdpcTx               = "ldpc_tx"
	IovarStrSgiRx                = "sgi_rx"
	IovarStrSgiTx                = "sgi_tx"
	IovarStrApivtwOverride       = "brcmapivtwo"
	IovarStrBwteBwteGciMask      = "bwte_gci_mask"
	IovarStrBwteGciSendmsg       = "bwte_gci_sendm"
	IovarStrWdDisable            = "wd_disable"
	IovarStrDltro                = "dltro"
	IovarStrSaePassword          = "sae_password"
	IovarStrSaePweLoop           = "sae_max_pwe_loop"
	IovarStrPmkidInfo            = "pmkid_info"
	IovarStrPmkidClear           = "pmkid_clear"
	IovarStrAuthStatus           = "auth_status"
	IovarStrBtcLescanParams      = "btc_lescan_params"
	IovarStrArpVersion           = "arp_version"
	IovarStrArpPeerage           = "arp_peerage"
	IovarStrArpoe                = "arpoe"
	IovarStrArpOl                = "arp_ol"
	IovarStrArpTableClear        = "arp_table_clear"
	IovarStrArpHostip            = "arp_hostip"
	IovarStrArpHostipClear       = "arp_hostip_clear"
	IovarStrArpStats             = "arp_stats"
	IovarStrArpStatsClear        = "arp_stats_clear"
	IovarStrTko                  = "tko"
	IovarStrRoamTimeThresh       = "roam_time_thresh"
	iovarWnmMaxidle              = "wnm_maxidle"
	IovarStrHe                   = "he"
	IovarStrTwt                  = "twt"
	IovarStrOffloadConfig        = "offload_config"
	IovarStrWsecInfo             = "wsec_info"
	IovarStrKeepaliveConfig      = "keep_alive"
	IovarStrMbo                  = "mbo"
)

const (
	wlcGetMagic                   = 0
	wlcGetVersion                 = 1
	wlcUp                         = 2
	wlcDown                       = 3
	wlcGetLoop                    = 4
	wlcSetLoop                    = 5
	wlcDump                       = 6
	wlcGetMsglevel                = 7
	wlcSetMsglevel                = 8
	wlcGetPromisc                 = 9
	wlcSetPromisc                 = 10
	wlcGetRate                    = 12
	wlcGetInstance                = 14
	wlcGetInfra                   = 19
	wlcSetInfra                   = 20
	wlcGetAuth                    = 21
	wlcSetAuth                    = 22
	wlcGetBssid                   = 23
	wlcSetBssid                   = 24
	wlcGetSsid                    = 25
	wlcSetSsid                    = 26
	wlcRestart                    = 27
	wlcGetChannel                 = 29
	wlcSetChannel                 = 30
	wlcGetSrl                     = 31
	wlcSetSrl                     = 32
	wlcGetLrl                     = 33
	wlcSetLrl                     = 34
	wlcGetPlcphdr                 = 35
	wlcSetPlcphdr                 = 36
	wlcGetRadio                   = 37
	wlcSetRadio                   = 38
	wlcGetPhytype                 = 39
	wlcDumpRate                   = 40
	wlcSetRateParams              = 41
	wlcGetKey                     = 44
	wlcSetKey                     = 45
	wlcGetRegulatory              = 46
	wlcSetRegulatory              = 47
	wlcGetPassiveScan             = 48
	wlcSetPassiveScan             = 49
	wlcScan                       = 50
	wlcScanResults                = 51
	wlcDisassoc                   = 52
	wlcReassoc                    = 53
	wlcGetRoamTrigger             = 54
	wlcSetRoamTrigger             = 55
	wlcGetRoamDelta               = 56
	wlcSetRoamDelta               = 57
	wlcGetRoamScanPeriod          = 58
	wlcSetRoamScanPeriod          = 59
	wlcEvm                        = 60
	wlcGetTxant                   = 61
	wlcSetTxant                   = 62
	wlcGetAntdiv                  = 63
	wlcSetAntdiv                  = 64
	wlcGetClosed                  = 67
	wlcSetClosed                  = 68
	wlcGetMaclist                 = 69
	wlcSetMaclist                 = 70
	wlcGetRateset                 = 71
	wlcSetRateset                 = 72
	wlcLongtrain                  = 74
	wlcGetBcnprd                  = 75
	wlcSetBcnprd                  = 76
	wlcGetDtimprd                 = 77
	wlcSetDtimprd                 = 78
	wlcGetSrom                    = 79
	wlcSetSrom                    = 80
	wlcGetWepRestrict             = 81
	wlcSetWepRestrict             = 82
	wlcGetCountry                 = 83
	wlcSetCountry                 = 84
	wlcGetPm                      = 85
	wlcSetPm                      = 86
	wlcGetWake                    = 87
	wlcSetWake                    = 88
	wlcGetForcelink               = 90
	wlcSetForcelink               = 91
	wlcFreqAccuracy               = 92
	wlcCarrierSuppress            = 93
	wlcGetPhyreg                  = 94
	wlcSetPhyreg                  = 95
	wlcGetRadioreg                = 96
	wlcSetRadioreg                = 97
	wlcGetRevinfo                 = 98
	wlcGetUcantdiv                = 99
	wlcSetUcantdiv                = 100
	wlcRReg                       = 101
	wlcWReg                       = 102
	wlcGetMacmode                 = 105
	wlcSetMacmode                 = 106
	wlcGetMonitor                 = 107
	wlcSetMonitor                 = 108
	wlcGetGmode                   = 109
	wlcSetGmode                   = 110
	wlcGetLegacyErp               = 111
	wlcSetLegacyErp               = 112
	wlcGetRxAnt                   = 113
	wlcGetCurrRateset             = 114
	wlcGetScansuppress            = 115
	wlcSetScansuppress            = 116
	wlcGetAp                      = 117
	wlcSetAp                      = 118
	wlcGetEapRestrict             = 119
	wlcSetEapRestrict             = 120
	wlcScbAuthorize               = 121
	wlcScbDeauthorize             = 122
	wlcGetWdslist                 = 123
	wlcSetWdslist                 = 124
	wlcGetAtim                    = 125
	wlcSetAtim                    = 126
	wlcGetRssi                    = 127
	wlcGetPhyantdiv               = 128
	wlcSetPhyantdiv               = 129
	wlcApRxOnly                   = 130
	wlcGetTxPathPwr               = 131
	wlcSetTxPathPwr               = 132
	wlcGetWsec                    = 133
	wlcSetWsec                    = 134
	wlcGetPhyNoise                = 135
	wlcGetBssInfo                 = 136
	wlcGetPktcnts                 = 137
	wlcGetLazywds                 = 138
	wlcSetLazywds                 = 139
	wlcGetBandlist                = 140
	wlcGetBand                    = 141
	wlcSetBand                    = 142
	wlcScbDeauthenticate          = 143
	wlcGetShortslot               = 144
	wlcGetShortslotOverride       = 145
	wlcSetShortslotOverride       = 146
	wlcGetShortslotRestrict       = 147
	wlcSetShortslotRestrict       = 148
	wlcGetGmodeProtection         = 149
	wlcGetGmodeProtectionOverride = 150
	wlcSetGmodeProtectionOverride = 151
	wlcUpgrade                    = 152
	wlcGetIgnoreBcns              = 155
	wlcSetIgnoreBcns              = 156
	wlcGetScbTimeout              = 157
	wlcSetScbTimeout              = 158
	wlcGetAssoclist               = 159
	wlcGetClk                     = 160
	wlcSetClk                     = 161
	wlcGetUp                      = 162
	wlcOut                        = 163
	wlcGetWpaAuth                 = 164
	wlcSetWpaAuth                 = 165
	wlcGetUcflags                 = 166
	wlcSetUcflags                 = 167
	wlcGetPwridx                  = 168
	wlcSetPwridx                  = 169
	wlcGetTssi                    = 170
	wlcGetSupRatesetOverride      = 171
	wlcSetSupRatesetOverride      = 172
	wlcGetProtectionControl       = 178
	wlcSetProtectionControl       = 179
	wlcGetPhylist                 = 180
	wlcEncryptStrength            = 181
	wlcDecryptStatus              = 182
	wlcGetKeySeq                  = 183
	wlcGetScanChannelTime         = 184
	wlcSetScanChannelTime         = 185
	wlcGetScanUnassocTime         = 186
	wlcSetScanUnassocTime         = 187
	wlcGetScanHomeTime            = 188
	wlcSetScanHomeTime            = 189
	wlcGetScanNprobes             = 190
	wlcSetScanNprobes             = 191
	wlcGetPrbRespTimeout          = 192
	wlcSetPrbRespTimeout          = 193
	wlcGetAtten                   = 194
	wlcSetAtten                   = 195
	wlcGetShmem                   = 196
	wlcSetShmem                   = 197
	wlcSetWsecTest                = 200
	wlcScbDeauthenticateForReason = 201
	wlcTkipCountermeasures        = 202
	wlcGetPiomode                 = 203
	wlcSetPiomode                 = 204
	wlcSetAssocPrefer             = 205
	wlcGetAssocPrefer             = 206
	wlcSetRoamPrefer              = 207
	wlcGetRoamPrefer              = 208
	wlcSetLed                     = 209
	wlcGetLed                     = 210
	wlcGetInterferenceMode        = 211
	wlcSetInterferenceMode        = 212
	wlcGetChannelQa               = 213
	wlcStartChannelQa             = 214
	wlcGetChannelSel              = 215
	wlcStartChannelSel            = 216
	wlcGetValidChannels           = 217
	wlcGetFakefrag                = 218
	wlcSetFakefrag                = 219
	wlcGetPwroutPercentage        = 220
	wlcSetPwroutPercentage        = 221
	wlcSetBadFramePreempt         = 222
	wlcGetBadFramePreempt         = 223
	wlcSetLeapList                = 224
	wlcGetLeapList                = 225
	wlcGetCwmin                   = 226
	wlcSetCwmin                   = 227
	wlcGetCwmax                   = 228
	wlcSetCwmax                   = 229
	wlcGetWet                     = 230
	wlcSetWet                     = 231
	wlcGetPub                     = 232
	wlcGetKeyPrimary              = 235
	wlcSetKeyPrimary              = 236
	wlcGetAciArgs                 = 238
	wlcSetAciArgs                 = 239
	wlcUnsetCallback              = 240
	wlcSetCallback                = 241
	wlcGetRadar                   = 242
	wlcSetRadar                   = 243
	wlcSetSpectManagment          = 244
	wlcGetSpectManagment          = 245
	wlcWdsGetRemoteHwaddr         = 246
	wlcWdsGetWpaSup               = 247
	wlcSetCsScanTimer             = 248
	wlcGetCsScanTimer             = 249
	wlcMeasureRequest             = 250
	wlcInit                       = 251
	wlcSendQuiet                  = 252
	wlcKeepalive                  = 253
	wlcSendPwrConstraint          = 254
	wlcUpgradeStatus              = 255
	wlcCurrentPwr                 = 256
	wlcGetScanPassiveTime         = 257
	wlcSetScanPassiveTime         = 258
	wlcLegacyLinkBehavior         = 259
	wlcGetChannelsInCountry       = 260
	wlcGetCountryList             = 261
	wlcGetVar                     = 262
	wlcSetVar                     = 263
	wlcNvramGet                   = 264
	wlcNvramSet                   = 265
	wlcNvramDump                  = 266
	wlcReboot                     = 267
	wlcSetWsecPmk                 = 268
	wlcGetAuthMode                = 269
	wlcSetAuthMode                = 270
	wlcGetWakeentry               = 271
	wlcSetWakeentry               = 272
	wlcNdconfigItem               = 273
	wlcNvotpw                     = 274
	wlcOtpw                       = 275
	wlcIovBlockGet                = 276
	wlcIovModulesGet              = 277
	wlcSoftReset                  = 278
	wlcGetAllowMode               = 279
	wlcSetAllowMode               = 280
	wlcGetDesiredBssid            = 281
	wlcSetDesiredBssid            = 282
	wlcDisassocMyap               = 283
	wlcGetNbands                  = 284
	wlcGetBandstates              = 285
	wlcGetWlcBssInfo              = 286
	wlcGetAssocInfo               = 287
	wlcGetOidPhy                  = 288
	wlcSetOidPhy                  = 289
	wlcSetAssocTime               = 290
	wlcGetDesiredSsid             = 291
	wlcGetChanspec                = 292
	wlcGetAssocState              = 293
	wlcSetPhyState                = 294
	wlcGetScanPending             = 295
	wlcGetScanreqPending          = 296
	wlcGetPrevRoamReason          = 297
	wlcSetPrevRoamReason          = 298
	wlcGetBandstatesPi            = 299
	wlcGetPhyState                = 300
	wlcGetBssWpaRsn               = 301
	wlcGetBssWpa2Rsn              = 302
	wlcGetBssBcnTs                = 303
	wlcGetIntDisassoc             = 304
	wlcSetNumPeers                = 305
	wlcGetNumBss                  = 306
	wlcGetWsecPmk                 = 318
	wlcGetRandomBytes             = 319
	wlcLast                       = 320
)

const (
	sdiodCccrRev          = 0x00  /* CCCR/SDIO Revision */
	sdiodCccrSdrev        = 0x01  /* SD Revision */
	sdiodCccrIoen         = 0x02  /* I/O Enable */
	sdiodCccrIordy        = 0x03  /* I/O Ready */
	sdiodCccrInten        = 0x04  /* Interrupt Enable */
	sdiodCccrIntpend      = 0x05  /* Interrupt Pending */
	sdiodCccrIoabort      = 0x06  /* I/O Abort */
	sdiodCccrBictrl       = 0x07  /* Bus Interface control */
	sdiodCccrCapablities  = 0x08  /* Card Capabilities */
	sdiodCccrCisptr0      = 0x09  /* Common CIS Base Address Pointer Register 0 (LSB) */
	sdiodCccrCisptr1      = 0x0A  /* Common CIS Base Address Pointer Register 1 */
	sdiodCccrCisptr2      = 0x0B  /* Common CIS Base Address Pointer Register 2 (MSB - only bit 1 valid)*/
	sdiodCccrBussusp      = 0x0C  /* Bus Suspend. Valid only if SBS is set */
	sdiodCccrFuncsel      = 0x0D  /* Function Select. Valid only if SBS is set */
	sdiodCccrExecflags    = 0x0E  /* Exec Flags. Valid only if SBS is set */
	sdiodCccrRdyflags     = 0x0F  /* Ready Flags. Valid only if SBS is set */
	sdiodCccrBlksize0     = 0x10  /* Function 0 (Bus) SDIO Block Size Register 0 (LSB) */
	sdiodCccrBlksize1     = 0x11  /* Function 0 (Bus) SDIO Block Size Register 1 (MSB) */
	sdiodCccrPowerControl = 0x12  /* Power Control */
	sdiodCccrSpeedControl = 0x13  /* Bus Speed Select  (control device entry into high-speed clocking mode)  */
	sdiodCccrUhsI         = 0x14  /* UHS-I Support */
	sdiodCccrDrive        = 0x15  /* Drive Strength */
	sdiodCccrIntext       = 0x16  /* Interrupt Extension */
	sdiodCccrBrcmCardcap  = 0xF0  /* Brcm Card Capability */
	sdiodCccrBrcmCardctl  = 0xF1  /* Brcm Card Control */
	sdiodSepIntCtl        = 0xF2  /* Separate Interrupt Control*/
	sdiodCccrF1info       = 0x100 /* Function 1 (Backplane) Info */
	sdiodCccrF1hp         = 0x102 /* Function 1 (Backplane) High Power */
	sdiodCccrF1cisptr0    = 0x109 /* Function 1 (Backplane) CIS Base Address Pointer Register 0 (LSB) */
	sdiodCccrF1cisptr1    = 0x10A /* Function 1 (Backplane) CIS Base Address Pointer Register 1       */
	sdiodCccrF1cisptr2    = 0x10B /* Function 1 (Backplane) CIS Base Address Pointer Register 2 (MSB - only bit 1 valid) */
	sdiodCccrF1blksize0   = 0x110 /* Function 1 (Backplane) SDIO Block Size Register 0 (LSB) */
	sdiodCccrF1blksize1   = 0x111 /* Function 1 (Backplane) SDIO Block Size Register 1 (MSB) */
	sdiodCccrF2info       = 0x200 /* Function 2 (WLAN Data FIFO) Info */
	sdiodCccrF2hp         = 0x202 /* Function 2 (WLAN Data FIFO) High Power */
	sdiodCccrF2cisptr0    = 0x209 /* Function 2 (WLAN Data FIFO) CIS Base Address Pointer Register 0 (LSB) */
	sdiodCccrF2cisptr1    = 0x20A /* Function 2 (WLAN Data FIFO) CIS Base Address Pointer Register 1       */
	sdiodCccrF2cisptr2    = 0x20B /* Function 2 (WLAN Data FIFO) CIS Base Address Pointer Register 2 (MSB - only bit 1 valid) */
	sdiodCccrF2blksize0   = 0x210 /* Function 2 (WLAN Data FIFO) SDIO Block Size Register 0 (LSB) */
	sdiodCccrF2blksize1   = 0x211 /* Function 2 (WLAN Data FIFO) SDIO Block Size Register 1 (MSB) */
	sdiodCccrF3info       = 0x300 /* Function 3 (Bluetooth Data FIFO) Info */
	sdiodCccrF3hp         = 0x302 /* Function 3 (Bluetooth Data FIFO) High Power */
	sdiodCccrF3cisptr0    = 0x309 /* Function 3 (Bluetooth Data FIFO) CIS Base Address Pointer Register 0 (LSB) */
	sdiodCccrF3cisptr1    = 0x30A /* Function 3 (Bluetooth Data FIFO) CIS Base Address Pointer Register 1       */
	sdiodCccrF3cisptr2    = 0x30B /* Function 3 (Bluetooth Data FIFO) CIS Base Address Pointer Register 2 (MSB - only bit 1 valid) */
	sdiodCccrF3blksize0   = 0x310 /* Function 3 (Bluetooth Data FIFO) SDIO Block Size Register 0 (LSB) */
	sdiodCccrF3blksize1   = 0x311 /* Function 3 (Bluetooth Data FIFO) SDIO Block Size Register 1 (MSB) */
)

/* SDIOD_CCCR_BRCM_CARDCAP Bits */
const (
	sdiodCccrBrcmCardcapCmd14Support  = 0x02 /* Supports CMD14 */
	sdiodCccrBrcmCardcapCmd14Ext      = 0x04 /* CMD14 is allowed in FSM command state */
	sdiodCccrBrcmCardcapCmdNodec      = 0x08 /* sdiod_aos does not decode any command */
	sdiodCccrBrcmCardcapSecureMode    = 0x80 /* Supports bootloader security */
	sdiodCccrBrcmCardcapChipidPresent = 0x40 /* Supports Chip ID Read from SDIO Core */
)

/* SDIO CORE CHIPID REGISTER */
const (
	sdioCoreChipidReg = 0x330
)

/* SDIO_CHIP_CLOCK_CSR Bits */
const (
	sbsdioForceAlp         = 0x01 /* Force ALP request to backplane */
	sbsdioForceHt          = 0x02 /* Force HT request to backplane */
	sbsdioForceIlp         = 0x04 /* Force ILP request to backplane */
	sbsdioAlpAvailReq      = 0x08 /* Make ALP ready (power up xtal) */
	sbsdioHtAvailReq       = 0x10 /* Make HT ready (power up PLL) */
	sbsdioForceHwClkreqOff = 0x20 /* Squelch clock requests from HW */
	sbsdioAlpAvail         = 0x40 /* Status: ALP is ready */
	sbsdioHtAvail          = 0x80 /* Status: HT is ready */
	sbsdioRev8HtAvail      = 0x40
	sbsdioRev8AlpAvail     = 0x80

	sbsdioFunc1Sbaddrlow  = 0x1000A /* SB Address Window Low (b15) */
	sbsdioFunc1Sbaddrmid  = 0x1000B /* SB Address Window Mid (b23:b16) */
	sbsdioFunc1Sbaddrhigh = 0x1000C /* SB Address Window High (b31:b24) */
	sbsdioDeviceCtl       = 0x10009 /* control busy signal generation */
	sbsdioDevctlAddrRst   = 0x40    /* Reset SB Address to default value */
)

/* SDIO_FRAME_CONTROL Bits */
const (
	sfcRfTerm   = 1 << 0 /* Read Frame Terminate */
	sfcWfTerm   = 1 << 1 /* Write Frame Terminate */
	sfcCrc4woos = 1 << 2 /* HW reports CRC error for write out of sync */
	sfcAbortall = 1 << 3 /* Abort cancels all in-progress frames */
)

/* SDIO_TO_SB_MAILBOX bits corresponding to intstatus bits */
const (
	smbNak    = 1 << 0 /* To SB Mailbox Frame NAK */
	smbIntAck = 1 << 1 /* To SB Mailbox Host Interrupt ACK */
	smbUseOob = 1 << 2 /* To SB Mailbox Use OOB Wakeup */
	smbDevInt = 1 << 3 /* To SB Mailbox Miscellaneous Interrupt */
)

/* SDIO_WAKEUP_CTRL bits */
const (
	sbsdioWctrlWlWakeTillAlpAvail = 1 << 0 /* WL_WakeTillAlpAvail bit */
	sbsdioWctrlWlWakeTillHtAvail  = 1 << 1 /* WL_WakeTillHTAvail bit */
)

/* SDIO_SLEEP_CSR bits */
const (
	sbsdioSlpcsrKeepWlKso = 1 << 0
	sbsdioSlpcsrWlDevon   = 1 << 1
)

/* To hostmail box data */
const (
	IHmbDataNakHandled = 0x0001 /* retransmit NAK'd frame */
	IHmbDataDevReady   = 0x0002 /* talk to host after enable */
	IHmbDataFc         = 0x0004 /* per prio flowcontrol update flag */
	IHmbDataFwReady    = 0x0008 /* fw ready for protocol activity */
	IHmbDataFwHalt     = 0x0010 /* firmware halted */
)

/* SDIOD_CCCR_IOEN Bits */
const (
	sdioFuncEnable1 = 0x02 /* function 1 I/O enable */
	sdioFuncEnable2 = 0x04 /* function 2 I/O enable */
	sdioFuncEnable3 = 0x08 /* function 3 I/O enable */
)

const (
	sdioGpioSelect           = 0x10005
	sdioGpioOutput           = 0x10006
	sdioGpioEnable           = 0x10007
	sdioFunction2Watermark   = 0x10008
	sdioDeviceControl        = 0x10009
	sdioBackplaneAddressLow  = 0x1000A
	sdioBackplaneAddressMid  = 0x1000B
	sdioBackplaneAddressHigh = 0x1000C
	sdioFrameControl         = 0x1000D
	sdioChipClockCsr         = 0x1000E
	sdioPullUp               = 0x1000F
	sdioReadFrameBcLow       = 0x1001B
	sdioReadFrameBcHigh      = 0x1001C
	sdioWakeupCtrl           = 0x1001E
	sdioSleepCsr             = 0x1001F
	iHmbSwMask               = 0x000000F0
)

const (
	iHmbFcState  = 1 << 4
	iHmbFcChange = 1 << 5
	iHmbFrameInd = 1 << 6
	iHmbHostInt  = 1 << 7
)

const (
	wlcSupStatusOffset     = 256
	wlcDot11ScStatusOffset = 512

	wlcEStatusSuccess     = 0  /** operation was successful */
	wlcEStatusFail        = 1  /** operation failed */
	wlcEStatusTimeout     = 2  /** operation timed out */
	wlcEStatusNoNetworks  = 3  /** failed due to no matching network found */
	wlcEStatusAbort       = 4  /** operation was aborted */
	wlcEStatusNoAck       = 5  /** protocol failure: packet not ack'd */
	wlcEStatusUnsolicited = 6  /** AUTH or ASSOC packet was unsolicited */
	wlcEStatusAttempt     = 7  /** attempt to assoc to an auto auth configuration */
	wlcEStatusPartial     = 8  /** scan results are incomplete */
	wlcEStatusNewscan     = 9  /** scan aborted by another scan */
	wlcEStatusNewassoc    = 10 /** scan aborted due to assoc in progress */
	wlcEStatus11hquiet    = 11 /** 802.11h quiet period started */
	wlcEStatusSuppress    = 12 /** user disabled scanning (WLC_SET_SCANSUPPRESS) */
	wlcEStatusNochans     = 13 /** no allowable channels to scan */
	wlcEStatusCcxfastrm   = 14 /** scan aborted due to CCX fast roam */
	wlcEStatusCsAbort     = 15 /** abort channel select */

	/* for WLC_SUP messages */
	wlcSupDisconnected   = 0 + wlcSupStatusOffset
	wlcSupConnecting     = 1 + wlcSupStatusOffset
	wlcSupIdrequired     = 2 + wlcSupStatusOffset
	wlcSupAuthenticating = 3 + wlcSupStatusOffset
	wlcSupAuthenticated  = 4 + wlcSupStatusOffset
	wlcSupKeyxchange     = 5 + wlcSupStatusOffset
	wlcSupKeyed          = 6 + wlcSupStatusOffset
	wlcSupTimeout        = 7 + wlcSupStatusOffset
	wlcSupLastBasicState = 8 + wlcSupStatusOffset

	/* Extended supplicant authentication states */
	wlcSupKeyxchangeWaitM1 = wlcSupAuthenticated     /** Waiting   to receive handshake msg M1 */
	wlcSupKeyxchangePrepM2 = wlcSupKeyxchange        /** Preparing to send    handshake msg M2 */
	wlcSupKeyxchangeWaitM3 = wlcSupLastBasicState    /** Waiting   to receive handshake msg M3 */
	wlcSupKeyxchangePrepM4 = 9 + wlcSupStatusOffset  /** Preparing to send    handshake msg M4 */
	wlcSupKeyxchangeWaitG1 = 10 + wlcSupStatusOffset /** Waiting   to receive handshake msg G1 */
	wlcSupKeyxchangePrepG2 = 11 + wlcSupStatusOffset /** Preparing to send    handshake msg G2 */

	wlcDot11ScSuccess                = 0 + wlcDot11ScStatusOffset  /* Successful */
	wlcDot11ScFailure                = 1 + wlcDot11ScStatusOffset  /* Unspecified failure */
	wlcDot11ScCapMismatch            = 10 + wlcDot11ScStatusOffset /* Cannot support all requested capabilities in the Capability Information field */
	wlcDot11ScReassocFail            = 11 + wlcDot11ScStatusOffset /* Reassociation denied due to inability to confirm that association exists */
	wlcDot11ScAssocFail              = 12 + wlcDot11ScStatusOffset /* Association denied due to reason outside the scope of this standard */
	wlcDot11ScAuthMismatch           = 13 + wlcDot11ScStatusOffset /* Responding station does not support the specified authentication algorithm */
	wlcDot11ScAuthSeq                = 14 + wlcDot11ScStatusOffset /* Received an Authentication frame with authentication transaction sequence number out of expected sequence */
	wlcDot11ScAuthChallengeFail      = 15 + wlcDot11ScStatusOffset /* Authentication rejected because of challenge failure */
	wlcDot11ScAuthTimeout            = 16 + wlcDot11ScStatusOffset /* Authentication rejected due to timeout waiting for next frame in sequence */
	wlcDot11ScAssocBusyFail          = 17 + wlcDot11ScStatusOffset /* Association denied because AP is unable to handle additional associated stations */
	wlcDot11ScAssocRateMismatch      = 18 + wlcDot11ScStatusOffset /* Association denied due to requesting station not supporting all of the data rates in the BSSBasicRateSet parameter */
	wlcDot11ScAssocShortRequired     = 19 + wlcDot11ScStatusOffset /* Association denied due to requesting station not supporting the Short Preamble option */
	wlcDot11ScAssocPbccRequired      = 20 + wlcDot11ScStatusOffset /* Association denied due to requesting  station not supporting the PBCC Modulation option */
	wlcDot11ScAssocAgilityRequired   = 21 + wlcDot11ScStatusOffset /* Association denied due to requesting station not supporting the Channel Agility option */
	wlcDot11ScAssocSpectrumRequired  = 22 + wlcDot11ScStatusOffset /* Association denied because Spectrum Management capability is required. */
	wlcDot11ScAssocBadPowerCap       = 23 + wlcDot11ScStatusOffset /* Association denied because the info in the Power Cap element is unacceptable. */
	wlcDot11ScAssocBadSupChannels    = 24 + wlcDot11ScStatusOffset /* Association denied because the info in the Supported Channel element is unacceptable */
	wlcDot11ScAssocShortslotRequired = 25 + wlcDot11ScStatusOffset /* Association denied due to requesting station not supporting the Short Slot Time option */
	wlcDot11ScAssocErpbccRequired    = 26 + wlcDot11ScStatusOffset /* Association denied due to requesting station not supporting the ER-PBCC Modulation option */
	wlcDot11ScAssocDssofdmRequired   = 27 + wlcDot11ScStatusOffset /* Association denied due to requesting station not supporting the DSS-OFDM option */
	wlcDot11ScDeclined               = 37 + wlcDot11ScStatusOffset /* request declined */
	wlcDot11ScInvalidParams          = 38 + wlcDot11ScStatusOffset /* One or more params have invalid values */
	wlcDot11ScInvalidAkmp            = 43 + wlcDot11ScStatusOffset /* Association denied due to invalid AKMP */
	wlcDot11ScInvalidMdid            = 54 + wlcDot11ScStatusOffset /* Association denied due to invalid MDID */
	wlcDot11ScInvalidFtie            = 55 + wlcDot11ScStatusOffset /* Association denied due to invalid FTIE */

	WlcEStatusForce32Bit = 0x7FFFFFFE /** Force enum to be stored in 32 bit variable */
)

const (
	wlcEPruneReasonOffset   = 256
	wlcESupReasonOffset     = 512
	wlcEDot11RcReasonOffset = 768

	/* roam reason codes */
	wlcEReasonInitialAssoc   = 0 /** initial assoc */
	wlcEReasonLowRssi        = 1 /** roamed due to low RSSI */
	wlcEReasonDeauth         = 2 /** roamed due to DEAUTH indication */
	wlcEReasonDisassoc       = 3 /** roamed due to DISASSOC indication */
	wlcEReasonBcnsLost       = 4 /** roamed due to lost beacons */
	wlcEReasonFastRoamFailed = 5 /** roamed due to fast roam failure */
	wlcEReasonDirectedRoam   = 6 /** roamed due to request by AP */
	wlcEReasonTspecRejected  = 7 /** roamed due to TSPEC rejection */
	wlcEReasonBetterAp       = 8 /** roamed due to finding better AP */

	/* NAN sub-events comes as a reason code with event as WLC_E_NAN */
	wlcENanEventStatusChg = 9  /* generated on any change in nan_mac status */
	wlcENanEventMerge     = 10 /* Merged to a NAN cluster */
	wlcENanEventStop      = 11 /* NAN stopped */
	wlcENanEventP2p       = 12 /* NAN P2P EVENT */

	/* XXX: Dont use below four events: They will be cleanup use WL_NAN_EVENT_POST_DISC */
	wlcENanEventWindowBeginP2p     = 13 /* Event for begin of P2P further availability window */
	wlcENanEventWindowBeginMesh    = 14
	wlcENanEventWindowBeginIbss    = 15
	wlcENanEventWindowBeginRanging = 16
	wlcENanEventPostDisc           = 17 /* Event for post discovery data */
	wlcENanEventDataIfAdd          = 18 /* Event for Data IF add */
	wlcENanEventDataPeerAdd        = 19 /* Event for peer add */

	/* nan 2.0 */
	wlcENanEventDataInd  = 20 /* Data Indication to Host */
	wlcENanEventDataConf = 21 /* Data Response to Host */
	wlcENanEventSdfRx    = 22 /* entire service discovery frame */
	wlcENanEventDataEnd  = 23
	wlcENanEventBcnRx    = 24 /* received beacon payload */

	/* prune reason codes */
	wlcEPruneEncrMismatch  = 1 + wlcEPruneReasonOffset  /** encryption mismatch */
	wlcEPruneBcastBssid    = 2 + wlcEPruneReasonOffset  /** AP uses a broadcast BSSID */
	wlcEPruneMacDeny       = 3 + wlcEPruneReasonOffset  /** STA's MAC addr is in AP's MAC deny list */
	wlcEPruneMacNa         = 4 + wlcEPruneReasonOffset  /** STA's MAC addr is not in AP's MAC allow list */
	wlcEPruneRegPassv      = 5 + wlcEPruneReasonOffset  /** AP not allowed due to regulatory restriction */
	wlcEPruneSpctMgmt      = 6 + wlcEPruneReasonOffset  /** AP does not support STA locale spectrum mgmt */
	wlcEPruneRadar         = 7 + wlcEPruneReasonOffset  /** AP is on a radar channel of STA locale */
	wlcERsnMismatch        = 8 + wlcEPruneReasonOffset  /** STA does not support AP's RSN */
	wlcEPruneNoCommonRates = 9 + wlcEPruneReasonOffset  /** No rates in common with AP */
	wlcEPruneBasicRates    = 10 + wlcEPruneReasonOffset /** STA does not support all basic rates of BSS */
	wlcEPruneCcxfastPrevap = 11 + wlcEPruneReasonOffset /** CCX FAST ROAM: prune previous AP */
	wlcEPruneCipherNa      = 12 + wlcEPruneReasonOffset /** BSS's cipher not supported */
	wlcEPruneKnownSta      = 13 + wlcEPruneReasonOffset /** AP is already known to us as a STA */
	wlcEPruneCcxfastDroam  = 14 + wlcEPruneReasonOffset /** CCX FAST ROAM: prune unqualified AP */
	wlcEPruneWdsPeer       = 15 + wlcEPruneReasonOffset /** AP is already known to us as a WDS peer */
	wlcEPruneQbssLoad      = 16 + wlcEPruneReasonOffset /** QBSS LOAD - AAC is too low */
	wlcEPruneHomeAp        = 17 + wlcEPruneReasonOffset /** prune home AP */
	wlcEPruneApBlocked     = 18 + wlcEPruneReasonOffset /** prune blocked AP */
	wlcEPruneNoDiagSupport = 19 + wlcEPruneReasonOffset /** prune due to diagnostic mode not supported */

	/* WPA failure reason codes carried in the WLC_E_PSK_SUP event */
	wlcESupOther          = 0 + wlcESupReasonOffset  /** Other reason */
	wlcESupDecryptKeyData = 1 + wlcESupReasonOffset  /** Decryption of key data failed */
	wlcESupBadUcastWep128 = 2 + wlcESupReasonOffset  /** Illegal use of ucast WEP128 */
	wlcESupBadUcastWep40  = 3 + wlcESupReasonOffset  /** Illegal use of ucast WEP40 */
	wlcESupUnsupKeyLen    = 4 + wlcESupReasonOffset  /** Unsupported key length */
	wlcESupPwKeyCipher    = 5 + wlcESupReasonOffset  /** Unicast cipher mismatch in pairwise key */
	wlcESupMsg3TooManyIe  = 6 + wlcESupReasonOffset  /** WPA IE contains > 1 RSN IE in key msg 3 */
	wlcESupMsg3IeMismatch = 7 + wlcESupReasonOffset  /** WPA IE mismatch in key message 3 */
	wlcESupNoInstallFlag  = 8 + wlcESupReasonOffset  /** INSTALL flag unset in 4-way msg */
	wlcESupMsg3NoGtk      = 9 + wlcESupReasonOffset  /** encapsulated GTK missing from msg 3 */
	wlcESupGrpKeyCipher   = 10 + wlcESupReasonOffset /** Multicast cipher mismatch in group key */
	wlcESupGrpMsg1NoGtk   = 11 + wlcESupReasonOffset /** encapsulated GTK missing from group msg 1 */
	wlcESupGtkDecryptFail = 12 + wlcESupReasonOffset /** GTK decrypt failure */
	wlcESupSendFail       = 13 + wlcESupReasonOffset /** message send failure */
	wlcESupDeauth         = 14 + wlcESupReasonOffset /** received FC_DEAUTH */
	wlcESupWpaPskTmo      = 15 + wlcESupReasonOffset /** WPA PSK 4-way handshake timeout */

	dot11RcReserved        = 0 + wlcEDot11RcReasonOffset  /* d11 RC reserved */
	dot11RcUnspecified     = 1 + wlcEDot11RcReasonOffset  /* Unspecified reason */
	dot11RcAuthInval       = 2 + wlcEDot11RcReasonOffset  /* Previous authentication no longer valid */
	dot11RcDeauthLeaving   = 3 + wlcEDot11RcReasonOffset  /* Deauthenticated because sending station is leaving (or has left) IBSS or ESS */
	dot11RcInactivity      = 4 + wlcEDot11RcReasonOffset  /* Disassociated due to inactivity */
	dot11RcBusy            = 5 + wlcEDot11RcReasonOffset  /* Disassociated because AP is unable to handle all currently associated stations */
	dot11RcInvalClass2     = 6 + wlcEDot11RcReasonOffset  /* Class 2 frame received from nonauthenticated station */
	dot11RcInvalClass3     = 7 + wlcEDot11RcReasonOffset  /* Class 3 frame received from nonassociated station */
	dot11RcDisassocLeaving = 8 + wlcEDot11RcReasonOffset  /* Disassociated because sending station is leaving (or has left) BSS */
	dot11RcNotAuth         = 9 + wlcEDot11RcReasonOffset  /* Station requesting (re)association is not * authenticated with responding station */
	dot11RcBadPc           = 10 + wlcEDot11RcReasonOffset /* Unacceptable power capability element */
	dot11RcBadChannels     = 11 + wlcEDot11RcReasonOffset /* Unacceptable supported channels element */
	/* 12 is unused */
	/* XXX 13-23 are WPA/802.11i reason codes defined in proto/wpa.h */
	/* 32-39 are QSTA specific reasons added in 11e */
	dot11RcUnspecifiedQos  = 32 + wlcEDot11RcReasonOffset /* unspecified QoS-related reason */
	dot11RcInsuffcientBw   = 33 + wlcEDot11RcReasonOffset /* QAP lacks sufficient bandwidth */
	dot11RcExcessiveFrames = 34 + wlcEDot11RcReasonOffset /* excessive number of frames need ack */
	dot11RcTxOutsideTxop   = 35 + wlcEDot11RcReasonOffset /* transmitting outside the limits of txop */
	dot11RcLeavingQbss     = 36 + wlcEDot11RcReasonOffset /* QSTA is leaving the QBSS (or restting) */
	dot11RcBadMechanism    = 37 + wlcEDot11RcReasonOffset /* does not want to use the mechanism */
	dot11RcSetupNeeded     = 38 + wlcEDot11RcReasonOffset /* mechanism needs a setup */
	dot11RcTimeout         = 39 + wlcEDot11RcReasonOffset /* timeout */
	dot11RcMax             = 23 + wlcEDot11RcReasonOffset /* Reason codes > 23 are reserved */

	wlcEReasonForce32Bit = 0x7FFFFFFE /** Force enum to be stored in 32 bit variable */
)

const (
	wlcENone                       = 0x7FFFFFFE
	wlcESetSsid                    = 0  /** indicates status of set SSID */
	wlcEJoin                       = 1  /** differentiates join IBSS from found (WLC_E_START) IBSS */
	wlcEStart                      = 2  /** STA founded an IBSS or AP started a BSS */
	wlcEAuth                       = 3  /** 802.11 AUTH request */
	wlcEAuthInd                    = 4  /** 802.11 AUTH indication */
	wlcEDeauth                     = 5  /** 802.11 DEAUTH request */
	wlcEDeauthInd                  = 6  /** 802.11 DEAUTH indication */
	wlcEAssoc                      = 7  /** 802.11 ASSOC request */
	wlcEAssocInd                   = 8  /** 802.11 ASSOC indication */
	wlcEReassoc                    = 9  /** 802.11 REASSOC request */
	wlcEReassocInd                 = 10 /** 802.11 REASSOC indication */
	wlcEDisassoc                   = 11 /** 802.11 DISASSOC request */
	wlcEDisassocInd                = 12 /** 802.11 DISASSOC indication */
	wlcEQuietStart                 = 13 /** 802.11h Quiet period started */
	wlcEQuietEnd                   = 14 /** 802.11h Quiet period ended */
	wlcEBeaconRx                   = 15 /** BEACONS received/lost indication */
	wlcELink                       = 16 /** generic link indication */
	wlcEMicError                   = 17 /** TKIP MIC error occurred */
	wlcENdisLink                   = 18 /** NDIS style link indication */
	wlcERoam                       = 19 /** roam attempt occurred: indicate status & reason */
	wlcETxfail                     = 20 /** change in dot11FailedCount (txfail) */
	wlcEPmkidCache                 = 21 /** WPA2 pmkid cache indication */
	wlcERetrogradeTsf              = 22 /** current AP's TSF value went backward */
	wlcEPrune                      = 23 /** AP was pruned from join list for reason */
	wlcEAutoauth                   = 24 /** report AutoAuth table entry match for join attempt */
	wlcEEapolMsg                   = 25 /** Event encapsulating an EAPOL message */
	wlcEScanComplete               = 26 /** Scan results are ready or scan was aborted */
	wlcEAddtsInd                   = 27 /** indicate to host addts fail/success */
	wlcEDeltsInd                   = 28 /** indicate to host delts fail/success */
	wlcEBcnsentInd                 = 29 /** indicate to host of beacon transmit */
	wlcEBcnrxMsg                   = 30 /** Send the received beacon up to the host */
	wlcEBcnlostMsg                 = 31 /** indicate to host loss of beacon */
	wlcERoamPrep                   = 32 /** before attempting to roam */
	wlcEPfnNetFound                = 33 /** PFN network found event */
	wlcEPfnNetLost                 = 34 /** PFN network lost event */
	wlcEResetComplete              = 35
	wlcEJoinStart                  = 36
	wlcERoamStart                  = 37
	wlcEAssocStart                 = 38
	wlcEIbssAssoc                  = 39
	wlcERadio                      = 40
	wlcEPsmWatchdog                = 41 /** PSM microcode watchdog fired */
	wlcECcxAssocStart              = 42 /** CCX association start */
	wlcECcxAssocAbort              = 43 /** CCX association abort */
	wlcEProbreqMsg                 = 44 /** probe request received */
	wlcEScanConfirmInd             = 45
	wlcEPskSup                     = 46 /** WPA Handshake */
	wlcECountryCodeChanged         = 47
	wlcEExceededMediumTime         = 48 /** WMMAC excedded medium time */
	wlcEIcvError                   = 49 /** WEP ICV error occurred */
	wlcEUnicastDecodeError         = 50 /** Unsupported unicast encrypted frame */
	wlcEMulticastDecodeError       = 51 /** Unsupported multicast encrypted frame */
	wlcETrace                      = 52
	wlcEBtaHciEvent                = 53 /** BT-AMP HCI event */
	wlcEIf                         = 54 /** I/F change (for wlan host notification) */
	wlcEP2pDiscListenComplete      = 55 /** P2P Discovery listen state expires */
	wlcERssi                       = 56 /** indicate RSSI change based on configured levels */
	wlcEPfnBestBatching            = 57 /** PFN best network batching event */
	wlcEExtlogMsg                  = 58
	wlcEActionFrame                = 59 /** Action frame reception */
	wlcEActionFrameComplete        = 60 /** Action frame Tx complete */
	wlcEPreAssocInd                = 61 /** assoc request received */
	wlcEPreReassocInd              = 62 /** re-assoc request received */
	wlcEChannelAdopted             = 63 /** channel adopted (xxx: obsoleted) */
	wlcEApStarted                  = 64 /** AP started */
	wlcEDfsApStop                  = 65 /** AP stopped due to DFS */
	wlcEDfsApResume                = 66 /** AP resumed due to DFS */
	wlcEWaiStaEvent                = 67 /** WAI stations event */
	wlcEWaiMsg                     = 68 /** event encapsulating an WAI message */
	wlcEEscanResult                = 69 /** escan result event */
	wlcEActionFrameOffChanComplete = 70 /** action frame off channel complete */ /* NOTE - This used to be WLC_E_WAKE_EVENT */
	wlcEProbrespMsg                = 71 /** probe response received */
	wlcEP2pProbreqMsg              = 72 /** P2P Probe request received */
	wlcEDcsRequest                 = 73
	wlcEFifoCreditMap              = 74 /** credits for D11 FIFOs. [AC0AC1AC2AC3BC_MCATIM] */
	wlcEActionFrameRx              = 75 /** Received action frame event WITH wl_event_rx_frame_data_t header */
	wlcEWakeEvent                  = 76 /** Wake Event timer fired used for wake WLAN test mode */
	wlcERmComplete                 = 77 /** Radio measurement complete */
	wlcEHtsfsync                   = 78 /** Synchronize TSF with the host */
	wlcEOverlayReq                 = 79 /** request an overlay IOCTL/iovar from the host */
	wlcECsaCompleteInd             = 80
	wlcEExcessPmWakeEvent          = 81 /** excess PM Wake Event to inform host  */
	wlcEPfnScanNone                = 82 /** no PFN networks around */
	wlcEPfnScanAllgone             = 83 /** last found PFN network gets lost */
	wlcEGtkPlumbed                 = 84
	wlcEAssocIndNdis               = 85 /** 802.11 ASSOC indication for NDIS only */
	wlcEReassocIndNdis             = 86 /** 802.11 REASSOC indication for NDIS only */
	wlcEAssocReqIe                 = 87
	wlcEAssocRespIe                = 88
	wlcEAssocRecreated             = 89  /** association recreated on resume */
	wlcEActionFrameRxNdis          = 90  /** rx action frame event for NDIS only */
	wlcEAuthReq                    = 91  /** authentication request received */
	wlcETdlsPeerEvent              = 92  /** discovered peer connected/disconnected peer */
	wlcEMeshDhcpSuccess            = 92  /** DHCP handshake successful for a mesh interface */
	wlcESpeedyRecreateFail         = 93  /** fast assoc recreation failed */
	wlcENative                     = 94  /** port-specific event and payload (e.g. NDIS) */
	wlcEPktdelayInd                = 95  /** event for tx pkt delay suddently jump */
	wlcEAwdlAw                     = 96  /** AWDL AW period starts */
	wlcEAwdlRole                   = 97  /** AWDL Master/Slave/NE master role event */
	wlcEAwdlEvent                  = 98  /** Generic AWDL event */
	wlcENicAfTxs                   = 99  /** NIC AF txstatus */
	wlcENan                        = 100 /** NAN event */
	wlcEBeaconFrameRx              = 101
	wlcEServiceFound               = 102 /** desired service found */
	wlcEGasFragmentRx              = 103 /** GAS fragment received */
	wlcEGasComplete                = 104 /** GAS sessions all complete */
	wlcEP2poAddDevice              = 105 /** New device found by p2p offload */
	wlcEP2poDelDevice              = 106 /** device has been removed by p2p offload */
	wlcEWnmStaSleep                = 107 /** WNM event to notify STA enter sleep mode */
	wlcETxfailThresh               = 108 /** Indication of MAC tx failures (exhaustion of 802.11 retries) exceeding threshold(s) */
	wlcEProxd                      = 109 /** Proximity Detection event */
	wlcEIbssCoalesce               = 110 /** IBSS Coalescing */
	wlcEMeshPaired                 = 110 /** Mesh peer found and paired */
	wlcEAwdlRxPrbResp              = 111 /** AWDL RX Probe response */
	wlcEAwdlRxActFrame             = 112 /** AWDL RX Action Frames */
	wlcEAwdlWowlNullpkt            = 113 /** AWDL Wowl nulls */
	wlcEAwdlPhycalStatus           = 114 /** AWDL Phycal status */
	wlcEAwdlOobAfStatus            = 115 /** AWDL OOB AF status */
	wlcEAwdlScanStatus             = 116 /** Interleaved Scan status */
	wlcEAwdlAwStart                = 117 /** AWDL AW Start */
	wlcEAwdlAwEnd                  = 118 /** AWDL AW End */
	wlcEAwdlAwExt                  = 119 /** AWDL AW Extensions */
	wlcEAwdlPeerCacheControl       = 120
	wlcECsaStartInd                = 121
	wlcECsaDoneInd                 = 122
	wlcECsaFailureInd              = 123
	wlcECcaChanQual                = 124 /** CCA based channel quality report */
	wlcEBssid                      = 125 /** to report change in BSSID while roaming */
	wlcETxStatError                = 126 /** tx error indication */
	wlcEBcmcCreditSupport          = 127 /** credit check for BCMC supported */
	wlcEPstaPrimaryIntfInd         = 128 /** psta primary interface indication */
	wlcEBtWifiHandoverReq          = 130 /* Handover Request Initiated */
	wlcESpwTxinhibit               = 131 /* Southpaw TxInhibit notification */
	wlcEFbtAuthReqInd              = 132 /* FBT Authentication Request Indication */
	wlcERssiLqm                    = 133 /* Enhancement addition for WLC_E_RSSI */
	wlcEPfnGscanFullResult         = 134 /* Full probe/beacon (IEs etc) results */
	wlcEPfnSwc                     = 135 /* Significant change in rssi of bssids being tracked */
	wlcEAuthorized                 = 136 /* a STA been authroized for traffic */
	wlcEProbreqMsgRx               = 137 /* probe req with wl_event_rx_frame_data_t header */
	wlcEPfnScanComplete            = 138 /* PFN completed scan of network list */
	wlcERmcEvent                   = 139 /* RMC Event */
	wlcEDpstaIntfInd               = 140 /* DPSTA interface indication */
	wlcERrm                        = 141 /* RRM Event */
	wlcEUlp                        = 146 /* ULP entry event */
	wlcETko                        = 151 /* TCP Keep Alive Offload Event */
	wlcEExtAuthReq                 = 187 /* authentication request received */
	wlcEExtAuthFrameRx             = 188 /* authentication request received */
	wlcEMgmtFrameTxstatus          = 189 /* mgmt frame Tx complete */
	wlcECsiEnable                  = 198 /* Setup a communication with the application layer to send CSI data */
	wlcECsiData                    = 199 /* Send the CSI data to application layer */
	wlcECsiDisable                 = 200 /* Diable the communication from application layer */
	wlcELast                       = 201 /* highest val + 1 for range checking */
)

type bssType int32

const (
	whdBssTypeInfrastructure bssType = 0  /**< Denotes infrastructure network                  */
	whdBssTypeAdhoc          bssType = 1  /**< Denotes an 802.11 ad-hoc IBSS network           */
	whdBssTypeAny            bssType = 2  /**< Denotes either infrastructure or ad-hoc network */
	whdBssTypeMesh           bssType = 3  /**< Denotes 802.11 mesh network                     */
	whdBssTypeUnknown        bssType = -1 /**< May be returned by scan function if BSS type is unknown. Do not pass this to the Join function */
)

type securityType uint32

const (
	wepEnabled    = 0x0001     /**< Flag to enable WEP Security        */
	tkipEnabled   = 0x0002     /**< Flag to enable TKIP Encryption     */
	aesEnabled    = 0x0004     /**< Flag to enable AES Encryption      */
	sharedEnabled = 0x00008000 /**< Flag to enable Shared key Security */
	wpaSecurity   = 0x00200000 /**< Flag to enable WPA Security        */
	wpa2Security  = 0x00400000 /**< Flag to enable WPA2 Security       */
	wpa3Security  = 0x01000000 /**< Flag to enable WPA3 PSK Security   */
	wpa3Owe       = 0x80000000 /**< Flag to enable WPA3 OWE Security   */

	enterpriseEnabled = 0x02000000 /**< Flag to enable Enterprise Security */
	sha2561x          = 0x04000000 /**< Flag 1X with SHA256 key derivation */
	suiteBSha384      = 0x08000000 /**< Flag to enable Suite B-192 SHA384 Security */
	wpsEnabled        = 0x10000000 /**< Flag to enable WPS Security        */
	ibssEnabled       = 0x20000000 /**< Flag to enable IBSS mode           */
	fbtEnabled        = 0x40000000 /**< Flag to enable FBT                 */

	whd_SECURITY_OPEN                securityType = 0                                                     /**< Open security                                         */
	whd_SECURITY_WEP_PSK             securityType = wepEnabled                                            /**< WEP PSK Security with open authentication             */
	whd_SECURITY_WEP_SHARED          securityType = wepEnabled | sharedEnabled                            /**< WEP PSK Security with shared authentication           */
	whd_SECURITY_WPA_TKIP_PSK        securityType = wpaSecurity | tkipEnabled                             /**< WPA PSK Security with TKIP                            */
	whd_SECURITY_WPA_AES_PSK         securityType = wpaSecurity | aesEnabled                              /**< WPA PSK Security with AES                             */
	whd_SECURITY_WPA_MIXED_PSK       securityType = wpaSecurity | aesEnabled | tkipEnabled                /**< WPA PSK Security with AES & TKIP                      */
	whd_SECURITY_WPA2_AES_PSK        securityType = wpa2Security | aesEnabled                             /**< WPA2 PSK Security with AES                            */
	whd_SECURITY_WPA2_AES_PSK_SHA256 securityType = wpa2Security | sha2561x | aesEnabled                  /**< WPA2 PSK SHA256 Security with AES                     */
	whd_SECURITY_WPA2_TKIP_PSK       securityType = wpa2Security | tkipEnabled                            /**< WPA2 PSK Security with TKIP                           */
	whd_SECURITY_WPA2_MIXED_PSK      securityType = wpa2Security | aesEnabled | tkipEnabled               /**< WPA2 PSK Security with AES & TKIP                     */
	whd_SECURITY_WPA2_FBT_PSK        securityType = wpa2Security | aesEnabled | fbtEnabled                /**< WPA2 FBT PSK Security with AES & TKIP */
	whd_SECURITY_WPA3_SAE            securityType = wpa3Security | aesEnabled                             /**< WPA3 Security with AES */
	whd_SECURITY_WPA3_FBT            securityType = wpa3Security | aesEnabled | fbtEnabled                /**< WPA3 Security with FBT                                */
	whd_SECURITY_WPA3_WPA2_PSK_FBT   securityType = wpa3Security | wpa2Security | aesEnabled | fbtEnabled /**< WPA3 WPA2 PSK security with AES FT enabled.           */
	whd_SECURITY_WPA2_WPA_AES_PSK    securityType = wpa2Security | wpaSecurity | aesEnabled               /**< WPA2 WPA PSK Security with AES                        */
	whd_SECURITY_WPA2_WPA_MIXED_PSK  securityType = wpa2Security | wpaSecurity | aesEnabled | tkipEnabled /**< WPA2 WPA PSK Security with AES & TKIP                  */
	whd_SECURITY_WPA3_WPA2_PSK       securityType = wpa3Security | wpa2Security | aesEnabled              /**< WPA3 WPA2 PSK Security with AES */

	whd_SECURITY_WPA_TKIP_ENT   securityType = enterpriseEnabled | wpaSecurity | tkipEnabled               /**< WPA Enterprise Security with TKIP                     */
	whd_SECURITY_WPA_AES_ENT    securityType = enterpriseEnabled | wpaSecurity | aesEnabled                /**< WPA Enterprise Security with AES                      */
	whd_SECURITY_WPA_MIXED_ENT  securityType = enterpriseEnabled | wpaSecurity | aesEnabled | tkipEnabled  /**< WPA Enterprise Security with AES & TKIP               */
	whd_SECURITY_WPA2_TKIP_ENT  securityType = enterpriseEnabled | wpa2Security | tkipEnabled              /**< WPA2 Enterprise Security with TKIP                    */
	whd_SECURITY_WPA2_AES_ENT   securityType = enterpriseEnabled | wpa2Security | aesEnabled               /**< WPA2 Enterprise Security with AES                     */
	whd_SECURITY_WPA2_MIXED_ENT securityType = enterpriseEnabled | wpa2Security | aesEnabled | tkipEnabled /**< WPA2 Enterprise Security with AES & TKIP              */
	whd_SECURITY_WPA2_FBT_ENT   securityType = enterpriseEnabled | wpa2Security | aesEnabled | fbtEnabled  /**< WPA2 Enterprise Security with AES & FBT               */

	whd_SECURITY_WPA3_192BIT_ENT   securityType = enterpriseEnabled | wpa3Security | suiteBSha384 | aesEnabled            /**< WPA3 192-BIT Enterprise Security with AES            */
	whd_SECURITY_WPA3_ENT          securityType = enterpriseEnabled | wpa3Security | sha2561x | aesEnabled                /**< WPA3 Enterprise Security with AES                    */
	whd_SECURITY_WPA3_ENT_AES_CCMP securityType = enterpriseEnabled | wpa3Security | wpa2Security | sha2561x | aesEnabled /**< WPA3 Enterprise transition Security with AES    */

	whd_SECURITY_IBSS_OPEN  securityType = ibssEnabled          /**< Open security on IBSS ad-hoc network                  */
	whd_SECURITY_WPS_SECURE securityType = aesEnabled           /**< WPS with AES security                                 */
	whd_SECURITY_wpa3Owe    securityType = wpa3Owe | aesEnabled /**< WPA3 Enhanced Open with AES security                  */

	whd_SECURITY_UNKNOWN securityType = 0xFFFF_FFFF // -1 /**< May be returned by scan function if security is unknown. Do not pass this to the join function! */

	whd_SECURITY_FORCE_32_BIT securityType = 0x7fffffff /**< Exists only to force whd_security_t type to 32 bits */
)

type band80211Type uint32

const (
	whd80211Band5ghz  band80211Type = 0 /**< Denotes 5GHz radio band   */
	whd80211Band24ghz band80211Type = 1 /**< Denotes 2.4GHz radio band */
	whd80211Band6ghz  band80211Type = 2 /**< Denotes 6GHz radio band   */
)

type scanType uint32

const (
	whdScanTypeActive             scanType = 0x00 /**< Actively scan a network by sending 802.11 probe(s)                              */
	whdScanTypePassive            scanType = 0x01 /**< Passively scan a network by listening for beacons from APs                      */
	whdScanTypePno                scanType = 0x02 /**< Use preferred network offload to detect an AP                                   */
	whdScanTypeProhibitedChannels scanType = 0x04 /**< Permit (passively) scanning a channel that isn't valid for the current country  */
	whdScanTypeNoBssidFilter      scanType = 0x08 /**< Return a scan record for each beacon or probe response RX'ed                    */
)

// WPA auth mode values for wlcSetWpaAuth IOCTL.
const (
	wpaAuthDisabled uint32 = 0x0000
	wpaAuthPsk      uint32 = 0x0004
	wpa2AuthPsk     uint32 = 0x0080
	wpa2AuthFt      uint32 = 0x1000
)

// wsec_pmk_t flags.
const (
	wsecPassphrase uint16 = 0x01
)

// wsec_pmk_t — passphrase/PMK structure for wlcSetWsecPmk.
type wsecPmkType struct {
	keyLen uint16
	flags  uint16
	key    [64]byte
}
