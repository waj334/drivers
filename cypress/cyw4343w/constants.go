package cyw4343w

const (
	busHeaderLen = 12
	ioctlOffset  = 4 + 12 + 16 // sizeof(void*) + 12 + 16
	ioctlMaxlen

	dataHeader       = 2
	asynceventHeader = 1
	controlHeader    = 0

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
