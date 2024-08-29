package main

import (
	"context"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"sort"
	"sync"
	"text/template"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	wg sync.WaitGroup
)

type TemplateData struct {
	Posts []*Post
}

type Post struct {
	Link      string
	Title     string
	Published time.Time
	Host      string
}

var (
	feeds = []string{
		"https://www.v2ex.com/feed/create.xml",
		"https://www.52pojie.cn/forum.php?mod=rss&fid=16",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCfbmmjCKPAxzLxeKqpg25vQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCX1ufnWuZ45jZeoJZQIUlFg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCwgFBJhSrb-wUTEeNkosqNw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCwVQIkAtyZzQSA-OY1rsGig",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCyEnglec9upW6pVDevuiCEQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC-jDxvKrM-QBqcDMzSvBn4A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCI9DPD1b5_y7ApDJ00jdi5Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCyWQTgUZuBGPhqQZYTizA3w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCPI2JV508nsMZrXQnMtRRZg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCK9RKx9b8X9M_BjSjJGDuyQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCjKfaVS0EIQJs_prvmJTejA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJzJdKsB2yrR4EJMN28cAqw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCHplVdsrnLDFar0Xd1C7pJA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCRUb1Iq2qbcmvt7AWdW8wjA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQMSVpxmoPD23fIQv4OCiZA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8L-NueEO99cfWT8IzTp1jQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1ISajIKfRN359MMmtckUTg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCEf_Bc-KVd7onSeifS3py9g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQi_oLFmCVZZN7yGtCQXx4g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSZQK43hJ422otti5D3iNgA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9xtexzjbEijc5Q6jGtwljQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCritGVo7pLJLUS8wEu32vow",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6jFCXt97O0_jhbKFVhwBGg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCCvt0Jc0ghFegppbyRdMPTg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCA_nEXFLmZyrgsxpjzYl6kA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCMUnInmOkrWN4gof9KlhNmQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC4zNtKX3XTplYQWNFPorvXQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCHBXStCavwkNt8XbvGbZtaw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1msa-mWNYZLeQfnwh510uQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCl1EDPkZuqwil6cgO_CSm8A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCIGilNMKaSuUHZclloJ5uIQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCmot9VXfYctSJSW_kNvoSaQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCpPLYBRk0Ax8uXJVtK5hw_w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCf6rOn9oohstgyGGpmP_hyQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCmHZ_X4N73fmD8lS3-Rgrzg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNE3Tho7XQ39iSUTafPmMpQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCr_F4Y9iboUKlg_ZPm4jkVQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCK8xfXnzhaARPLUh30uxtPw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCnNL5kFtW3ybCo7WSIj8cNg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCRhkaWPNfdGe4ZFdDDrLZAg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCIXhh33Q_2Cj1vMrsOIfXfA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9QYsNevKkWETe05fCt89Og",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCWIxaT-hfycZrtst0v3ITMw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCI1RTYd_EM0v5A-lVVwQgvw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOa8dTenJ_ylq2SaJNqpB9Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0JPCxTmOwY6G5L237d9R_g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCHVn2tGrSbjWVo578lXHUpA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCM3mfYn-PhhnvaIhK14dQSw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChAfs6xW5S6pfG7rllMPkTA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBAQY_7IdHjcBYfMzuVDAXA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCfTud6NK3SZ_k2BGnXNqxXg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC00ucsOFc-RsdGDOz0NHgSw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCFmFolFKjSIJceK7qFntDkg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCjAlA5tKCLb1QimeUIgdhqQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCaTCd7VJgJQXxloLnWuCkiA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCn14zhm8Zpx8IMr1zdHBoZg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCczA3DqHB2iV0gE5jNggq7w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCh21enuxt6axA0nkT6y4J2Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC_GcrznCXn6b-i1Y7DOn54A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJEER74X9kBenMT_x9iK9Mw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC5BMQOsAB8hKUyHu9KI6yig",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCtCiO5t2voB14CmZKTkIzPQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCqwUnggBBct-AY2lAdI88jQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9zY_E8mcAo_Oq772LEZq8Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCgD0APk2x9uBlLM0UsmhQjw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCeLPm9yH_a_QH8n6445G-Ow",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCqq-ovGE01ErlXakPihhKDA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChomMws0uvLe0HCN9H973Iw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCEIi7zFR_wE23jFncVtd6-A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSJ01G-E87y5ojdQKCXq7zg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCLYNEaV9stLOz2Wz04rqwhA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvQxfg4eOqzrpGb3awGM4Rg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCU0PgOXf0lxzVxN2TLzMJkw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQG8tNnV4hKetLhMb4MopHQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC5s_kUbxX3P1q6lmDgygD-w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCEz-AFAg3EUKsxraad1puQA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCzgxx_DM2Dcb9Y1spb9mUJA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCDNDlqJRz4FsO_ByfUNOSuQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6s5-dHFPVOSMH-Bw49gPAA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCZg7Xo4Ir4bB7Qdo7tXCHqg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCEpFoWeCMCo5z3EvWaz6hQQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCY3G7UI7_qdHRGH9HJXCTrA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCAJ6JgqsMkEFkiG82vxF1pA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCmot9VXfYctSJSW_kNvoSaQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCbU86dKACh3g1qiV3kMaRFQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCljnFhpcj6Ddic0-O-g3SwQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCgkQKV-Jum6D1hhx65jA24g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCPuXEvuQxa9Sj_9J1LaUB6g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCcnsf1dzzmkHg6SDl63L6eg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCTeZAwIRIbFgwJ8EacdO4dg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCxjXU89x6owat9dA8Z-bzdw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCTfgLRLltCx5ThH-3K8iAWw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCsBNDJamrYl4cIK4onhscPQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCbZiG8_Zmf1aOuw56cvX-Vw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0D9ASWn4HNR_nkpFT7Gzhg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCe7NidSvGNk8P3m5pI7oMVw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChiMMOhl6FpzjoRqvZ5rcaA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCgrw0YdNdFZnFuNZ96UWp_g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJDUoUiPR9ZwPZvwLpth7Lg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBkYT6mlOpVVV9t5nRHn9Ow",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1DBvdTqEkoXzza8UQsTFGQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC11HkuejItJOJm9BPnyBQtA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCpBhFFpFqoRxDoC557TygIA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9Tr6qvOJ99A2dM5QBQYwTQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCcMw8KMikSUKp1gcfYk0qsg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9K0rLE1SMh86nVxzkCBpNA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOkt-hnOi6T3ODq6pLDkgWw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCYjB6uufPeHSwuHs8wovLjg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCMlREL6GWSV2sfa8c0Uv82A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCCcK3JRvJWPvOfw2nKNmIIg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCXmPpvl1aVvSALP0EQ4Cexg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCV7nsHpgs0n8Q69eHXCZjLw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCEOhpXTO8CwrpkBTxOOu4bw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCICuPsiRtmuZ3T_EtLy7CRQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC46WXN1W1V9QYwGrjKRCiXg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC992aW0kfy4lttcg1ESSBWA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC84whx2xxsiA1gXHXXqKGOA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvSiZyvZvG5C89KcA3-qHnw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOxqgCwgOqC2lMqC5PYz_Dg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCRf1oaDTDdMDuI5NWoJN9tQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCKlqPlMAiORObF9rEGLXpgg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCHmDNBMqq4k7vaAbSZjjVFg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChmCw-lM24QMTg9BdA150tQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCILcVG4mRMoSs2meWA9qncA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvjB-wUbS7-HbIfDxjk2a3w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCp349XwMepHHFv64P__97XQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCofeVpB7pRjurJ6Aeejf3CQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQVR7lLqD720fIwpNFBvgbQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1WBWfOoJyuPkfObmex1AGA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkd3AEMr2zuV8eKeOT8G-NQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0BCfHbNORhd7dsw1R5rw9w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC7RNG2fClJFWXzwHb1Ym1DA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCzcEnsTFjLJWEpEKhcg9LAg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOFBs8PJIrnmi9N3WnDxKYg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCIA3-TFbWEk0EOtM4ZTyuHg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCGmip6pGWiI3iQbze350xdA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkwCj0YUftbLxqVx5fsblYg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC2nsP5D6eG2RSFLqA24AAsw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCrinMxQ_y9Iv0l_A9N4HnLw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC-Fnix71vRP64WXeo0ikd0Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0SkNQXPJ60hKEFubOz0fDA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC7lvKy6aBWRbgkk84uNoC5g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJ8A8IRA6Z0CeUqf9kCqGVA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC4GZ1dNQKWWFDQ4IWl4DezA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvBGKwlGtLjl7n4qmUfHiIg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCqFbF_CD2CJ_jwiQCfCJqEQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC_9IQjvjWRx3dDMzzgmYhTg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCLnFYJtpa_-rxW1KqkfJjbQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCV-dC5CIZBknF0lkoBduc_w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCHiof82PvgZrXFF-BRMvGDg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCxw_JBZaRWtL6JKvyxdUx5g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCFmL725KKPx2URVPvH3Gp8w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UClqNSqnWeOOUVkzcJFj4rBw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCdvYSTbhmzWgWyfGnhet03Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJyZfWrqaGX4nwXGKOEdM6Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCC2xZiqsdhsTd9AaC0achGQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCv5LdaIeRGjCYa0ppZH4fqA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCRz3cGfqeMPSHMBN6CxKQ9w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCxQiPuwF3iaxV_aJOzjQV-Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCb243laxn6xnsmtznXjgbnw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6VLBHsXtXtEMA3HpwAM2Rw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCjypzyxU_p0_0ElbWSi8QQw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC-MtmLd4PbmNuRaasDgBFBA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCIKOy_q2VWDv1vzeoi7KgNw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6gLlIAnzg7eJ8VuXDCZ_vg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCIwmqT5yaBLP4l5Emw7h35A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCdld7SHk9IyYSkNrHA9B6Dw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCS5qmPFj36fK7fML6gYcZSQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCp5Hfz5DiCueMcUzMvqU2eQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCDW9XFv3T8VkBskvgEIiW2A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCDvazhEPkMkIsDJGZsGw1zQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOON7HbktJy5i19fyKTOEXw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCzQ3t4sWb66begG8ZD11pdQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCHVRr6s2jTHYkcRgV_Z2rVw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCa3Qv6giqAn7qQTqm2ZnZXg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC4xgghZYLhpvrOIbSVH-NjQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9HHtF1CHSLxlQWXHjLVNVg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0Coab7B6mQRLs3eeaHBD7g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC5LH8st_ktFAV__AteRemzQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBzZXqzT7uDBD1_jeq2D0XQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNcPSDX4vQmNt-m8O7DEh5w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCem3aY5J_4PDC6-J-d7S0aQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCzRAXCeq4YaX6Y0ad7CwPEg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC2EP3T3-G0h-pP-Ovgnb07Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCVmw1jDjllTcAszVW3vWtag",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQuDYVdemGofgSX9LujILcQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UClVZ1wxpfCEF3ugZPkk3Zdw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCsBS3T8hYnflLJl-_0MpapA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCL1VUVgulGLH4idpAE60FFg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCblueD73Bc1wJqv3QRt0WUQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC2y0ntANyCqr7LxFCTsUoGQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCZgfwLNk27K0PdE8EGYNcPA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCC40X74K26Yp7SGxaPmq9Gg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCPITkokyrBnC-f40ocMo0XA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCps966LIRe60butjrVfSf8A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1opHUrw8rvnsadT-iGp7Cg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCVuGdRELvG83vYT7obfSqwQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvb-tSnkDvA05i6ujKyszjQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQERloooajfOt4WxauNaGRw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBo8Ygnvu36nyd2xlWD5amw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCiorb8zE67MvrBQGagNWutQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCLP31nVpk3jhhrbUuzbRbYQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCuX4Yhwl3ARiKe9iTxdCD4Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCY-toeMjzGcy7PVV0P8QcnA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC_ARaeDGHVLAqi6whEWQRTg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCnR2VoUqv9OF2oW4OmGZMzg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UClMJgjg2z_IrRm6J9KrhcuQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCUuENVpVuzqpRsXWIDlpQTg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC80UDR9O8vDvUVossfppxHg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCd0niLuvSNQnyOMlVuG-XyA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC5X7iRFTg24q8X01_W42OGg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCsr1UdXXySNlxgktuuiEt4g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNlMeUt5nOTQ-yfjXzRKVKA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCDtAY4AB4OgAZh9Sxovx2lA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvUmRTy-LMZYUrU6OoOJLMw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCAZpHV_W-0lA4BJuorbn-FQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCKiWll5qSb4oaEtZhYDtjyw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCZrK8gFUkbWdhk6J6I5wVkA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCWlL0DMfc2dlDhIx1kaGpKg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCoNfsDH8sZe13u7rSxaEBkw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCnEOcN_3q0M45fSqhFnN7UA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCdYXNCneIwpPBQdgRzV-ARA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1BIwK6cjVaadMUeYqmFDQA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCEk2wuevkGZFQ5BYrGgcFYw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCDc-k__VyoGLDgc-grT7Wcw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNdVOhYR2u-HiITV8oWEUBw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvUTCi2Ms1eP-bFnh0lfGHA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6zXBxeRnQOmyfRxtS3N8GQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNkX5VXtTe4Clkogm6h6QWA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCbE_Vqbd9SqiQFXnbm0Zrlw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6sMMhOdJu53G3doJGmsSXQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9XVPROR-gM3swX0EIOiRaw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCLnac6g3R8YcGOisZxsVIpw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCUj7ZoTL3RGXoFG8wyAF0Dg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8Ly-4dnxWOV9iddXO2WUdA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCWXwLgVslvBpzbGBz4lj4-Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCGb2oR63mwUexYDeHsnQwqQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1FCeHtSXC7kgQ9cf6bLfvg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0P0EBfu_VlIjrNc8PM2krg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCl1mdpP7n5LtXb81g-aYeIw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8_MZgNVeE5ff4uFX0zAbxg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCf1WHSQBcdOxh5Lxi0u-44w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCFMypDKGBtonh4yE7eKCzlA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQvsFc6XtUzKS0r9DNqxo3g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC_f62mghg-kP2WiY9MmhvPg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCN-PRV5vR_P6zEwfBXNkNbg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSWAuBR0Xk4m8XZjxuNqBzQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBhsePA4-EGWCJ_6tf3bODA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCCb9Uk8jTWP4XuvfOGrCq_Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJPgSXv_6ys7LZRJlb0lVRw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCe5nzJzWEYHOK938f-C5eMg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOcZpyUCrFchLeNoCPLusKg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCahfFeIvTMcqHiwsuBO-2Fw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCys-W1DKy0oGnJE7rJlbKUQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCgvJcttwBMvxXYiDfdvfubA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCiKoX0tQ7GArOoxnOiwch7g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCoxDu8GFUrENMIawkWkl37Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNoLaySBtTYWKjsNneyEbRg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCIKkjZ4_UUWnZxLIATgZvlQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBF9FTdOi_c414HylU8nG0w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC_z7H3yJ7u4usCa9O9dDJmw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCp7HCH-GO90GjQ0jx4jYVbg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8pByEpwbmRBy8O-OOub1aA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQ5RafKgkM1_WO4EfGTeKGg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCT-FiEMGRS_jK1CaFcCMzGw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCD3jscp2wLZyFG2X75b9paA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChpp6Q-sk-rhpt3fTOusFSw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCve-qnLI-_m2Oys2p3iW5Jw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvWqg1Xxcm6d0ltwOac9slA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCrdOI6gi0jAEhhg6kbVyRHQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQqpw2-neEiAHWsU5hv1xDw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCqZ7jCucPE9FeClgJhC3PPQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC3rrCl8CcCFDCI9Ti2rHsnw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCxGeii0KDrJVo-3GB1ur37Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6JBBjGhWrtkwOhREANselQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCh-ui5GFGcyxV8Sn-Tj557w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCt2rixQd2RIph9Ya2Ovij4g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCU3kq6pv8zS2K_xJ-j2M-KQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCCb6hI74fwNNMWV0cEPbXlw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC3WUnrStsLqG13pTnHECzqQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0zbE3G_hqpsEwHz4ig_Tww",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCir2AjQrwccbRSd_9o64GKQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCPpYonu8HMQCvNVcG2xYgwA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCHmmTQp2cnoCUsCLYHecsng",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCCkhMuD01VMVIxjgyZYRhFQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCht5-3fFIsGvPBS6UYH_Dnw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC3SyT4_WLHzN7JmHQwKQZww",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6jH57FvFkAl-7npkpMqWvg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCENYXixxLePOZY720oa6JfQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCYjd-3IqXi8WGxijxghJz1g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6LljJcM26j9A_2JCHM8FDw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCiAYMvLm2TcKlv7UZGZQqsQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCwBK59MFAsoEsMr_QTuag6A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCTs-yFOmlTSXQ50vmqksRMA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkrqHoLiY4EZFbWJAKO13kg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8-5FyWZDTVLosoNqcfD1ng",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBD69nZPJRjgJ1QxJ2TCHVw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCx6eBtnyfV5HDz2ePkuDLhQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCmbOr3Lf46Y9P9Lrwt1IDOQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC5oQTNq-gBURCIIgDzScT8w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCCjpp4TCjnRftS_LenXuSHw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCE6acMV3m35znLcf0JGNn7Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC86Z-Hn0tiGI5aXdx60DigQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkvg7ur8e_HBZTRDv7-RfjA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkOpWtMbYjN-JgSmmokWJ-Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCAtFkapSeoEGPxm5bC3tvaw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCo-gAYrvd7WIrCRsNueddtQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCDTHLADy7TuISpSlj3mQxfg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCCG5x9Hy7Ik5c0kpOPIpWfA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCuSxcjp1mVzMpGjMpLRpLxw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSG228wEbdvnIzjlkaFrGyw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCI7ktPB6toqucpkkCiolwLg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8k5iZIdewPjqwKtH7xk_yw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJXuQ9aRkS_HaFEBhiHFpWw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCTIkPi38UqlElB2m3N0erdw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCXe8hqIk_Yap31y7EOCMdzg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC5K0G6pTq2rSOhabuzBTHAA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCXY68I59Xq3dNG4z0L91H2Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCz4DsYIwsbL_hubAPNbUNuA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCb1oG-4QYwEJ_wLeXda-MhA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNRf6H1lLagyBirziAmcEbw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCW1ZwFNYFwgpKBsrZW-mPyg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNvQ9sosA89LGcnJODNbGsw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0max5gT8f5pKRl0BvhSlzw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCcNU8mZivckGgpeejAxYNXw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UClN5EsofXLP_SJae9WuzD0A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UClTsyw0gGmFrtjWIZCN1isg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCTyrENfBc1kF_WjJpJ69xwg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC5hrHnzmiahOE6vN8ZP5FKg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCf177kjSQKWZneInEnxjXpg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCdw86OcXaagJE0QQNa2NMyA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCw8hro4Xw7TiXQRNTnpV1Xw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBe1Ohvh25Q271AqYnktxDQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCmZb4LwQRhEzZX5Uqpcqziw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCEqHu6TEeyvTCRKQNl8dmoA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCQe2Y7V-C9bNMAcCJCBvzQQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCVpcg04wqPgtbUtsTVLmnAw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNFd_XQwGCsRTc2B8gD1s7A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCAHlZTSgcwNNpf8LV3E6kDQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCrTCkJ8377E7JhsMeyLXtKA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCadluGDMYnPhs4OtVpJXYoA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCaiF3ngpgU2Y-l9fjmLVe7w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNKkMtbbrSOEQWabCtyf9ZA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCgASMG4hbrxK543re59k86g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC_oO4HISr99F_vEP0RftrfQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCaCru8sowL6nlxkaoKAAcqw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSfgwPqPchZZK12QcUXDMHg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCp-ME1RUkd5y3dVlyo3Qvhg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC0jewycCEYd__JsHGCvIbFQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCzmHYNLBW2oSYZHRP75rAUw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvkWnAsnoF79YkhgssdEomA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCrx6P_kTFJguNkefyvA8a4w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCJJrtJnarZK_r7bncCdR5pQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC56FDco95hMXN33MWuUKRIQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCpYZtBFAs3roLhG0J30rxTw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8J54x9slo75kRvwf4Fp2Cw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCgeTlp-3-qe_pxiUVSif6hQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNe4AOBDQ9a5tV6ik0JuSVg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChs3VNBroEiMDwxyjVqnLKQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCNNxIP5s-YLQ4RfTrj5XOwA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC1gmZg_re0ptBDnFjbmWuhQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8kqQouREx6t8q5uGdKfRJA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC3EWyNiOL9Y0p0m_VSvDwwA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC6rJuYMZWNxG-GQ0LIQr3kg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCvtO27QvzLYfMwINGCl_gBg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCwSOARsvB-Qa6PtuYyK74dA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCbCCUH8S3yhlm7__rhxR2QQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8wI85SlecrsLdPm-I73PzA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC9ITbiOx3odg1aW2vEe1B-w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBhnI7_z3NIJmneNe7wjlPw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCKQArOKTvHlyp5aGaAXgJxQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCO2GVjlhscG3hxfOZYUYLxg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOo2BVhH0-f03z8IhtG3msw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCPL1pFqXQOqYeLe-PRYJYBg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSh7_4LCPpQabbsInWFv4Yg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCYSmtE_fGgruXOl1UoR52kQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCk3d6Oa_O8I7qogcxGeEm2w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCnEd5KvLOw7BTI4K5vEY1BA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC_9eLT-y0NkwlbcZaZKN7Mg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC00ucsOFc-RsdGDOz0NHgSw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCfYJ-ioTfZUhY59yBgojdTQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCAfIN6Cy_lgvFwoXXKQC2NA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCBe3AJEV4vXAH5xs7eVyTTA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCZUe9BnHZcqJI8yE-y-MKfQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkK2B6D3imy6EnpqfrYm-5A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkuisk2sEog8S4nMs6j3QUw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC-hTwjkGlNB5vyR6qccXr7A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCnY5pPAUxfPcxKEOjoPt-qg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCGwF8Vqxv6-1Fuhp7_cnu_w",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCOHAYSieHVvG8OEwanNVdLA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCB0XeIQhS8GrUA_LrzyyOEw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCh6snexfUrSdwKGuiyDSe5g",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCt8b5E7gf-OlfNasgjmHACw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC67o6aRi6mRsPh21Uw8YP2A",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCYM6U6usGTEl9gtxRzyMOiw",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCVTE4smKJV_27uQ682j3bIg",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCXuRgZDm5lPBNA8ddEDUrPQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCzYgQl2ML9JMC6zb8_ZIgrA",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCnr8XmfD3_2Hib3U0mGPWgQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCwAVxGVhFwNeJajMZ_cXa2Q",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCrvFGTw3NYCRyzYWAnf7UmQ",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCwXF2kyXFOCkPLZYC3sA_cQ",
		"https://www.ruanyifeng.com/blog/atom.xml"
	}

	// Show up to 60 days of posts
	relevantDuration = 5500 * 24 * time.Hour

	outputDir  = "docs" // So we can host the site on GitHub Pages
	outputFile = "index.html"

	// Error out if fetching feeds takes longer than a minute
	timeout = time.Minute
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	posts := getAllPosts(ctx, feeds)

	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return err
	}

	f, err := os.Create(path.Join(outputDir, outputFile))
	if err != nil {
		return err
	}
	defer f.Close()

	templateData := &TemplateData{
		Posts: posts,
	}

	if err := executeTemplate(f, templateData); err != nil {
		return err
	}

	return nil
}

// getAllPosts returns all posts from all feeds from the last `relevantDuration`
// time period. Posts are sorted chronologically descending.
func getAllPosts(ctx context.Context, feeds []string) []*Post {
	postChan := make(chan *Post)

	wg.Add(len(feeds))
	for _, feed := range feeds {
		go getPosts(ctx, feed, postChan)
	}

	var posts []*Post
	go func() {
		for post := range postChan {
			posts = append(posts, post)
		}
	}()

	wg.Wait()
	close(postChan)

	// Sort items chronologically descending
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Published.After(posts[j].Published)
	})

	return posts
}

func getPosts(ctx context.Context, feedURL string, posts chan *Post) {
	defer wg.Done()
	parser := gofeed.NewParser()
	feed, err := parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range feed.Items {
		published := item.PublishedParsed
		if published == nil {
			published = item.UpdatedParsed
		}
		if published.Before(time.Now().Add(-relevantDuration)) {
			continue
		}
		parsedLink, err := url.Parse(item.Link)
		if err != nil {
			log.Println(err)
		}
		post := &Post{
			Link:      item.Link,
			Title:     item.Title,
			Published: *published,
			Host:      parsedLink.Host,
		}
		posts <- post
	}
}

func executeTemplate(writer io.Writer, templateData *TemplateData) error {
	htmlTemplate := `
<!DOCTYPE html>
<html>
	<head>
	<link rel="icon" type="image/ico" href="https://jsd.onmicrosoft.cn/gh/rcy1314/tuchuang@main/NV/Level_Up_Your_Faith!_-_Geeks_Under_Grace.1yc7qyib5tsw.png">
    <link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/4.4.1/css/bootstrap.min.css">
    <link rel="stylesheet" href="https://cdn.staticfile.org/font-awesome/5.12.1/css/all.min.css">
	<link rel="stylesheet" href="ind.css">
    <link rel="stylesheet" href="style.css">
    <link rel="stylesheet" href="APlayer.min.css">
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>NOISE | 聚合信息阅读</title>
		<style>
		@import url("https://fonts.googleapis.com/css2?family=Nanum+Myeongjo&display=swap");

		body {
			font-family: "Nanum Myeongjo", serif;
			line-height: 1.7;
			max-width: 800px;
			margin:  auto ;
			padding: auto;
			height: 100%;
		}

		li {
			padding-bottom: 16px;
		}
	</style>
	</head>
	<script type='text/javascript' src='js/jquery-3.2.1.js'></script>  
        <script type='text/javascript'>  
            //显隐按钮  
            function showReposBtn(){  
                var clientHeight = $(window).height();  
                var scrollTop = $(document).scrollTop();  
                var maxScroll = $(document).height() - clientHeight;  
                //滚动距离超过可视一屏的距离时显示返回顶部按钮  
                if( scrollTop > clientHeight ){  
                    $('#retopbtn').show();  
                }else{  
                    $('#retopbtn').hide();  
                }  
                //滚动距离到达最底部时隐藏返回底部按钮  
                if( scrollTop >= maxScroll ){  
                    $('#rebtmbtn').hide();  
                }else{  
                    $('#rebtmbtn').show();  
                }  
            }  
              
            window.onload = function(){  
                //获取文档对象  
                $body = (window.opera) ? (document.compatMode == "CSS1Compat" ? $("html") : $("body")) : $("html,body");  
                //显示按钮  
                showReposBtn();  
            }  
              
            window.onscroll = function(){  
                //滚动时调整按钮显隐  
                showReposBtn();  
            }  
              
            //返回顶部  
            function returnTop(){  
                $body.animate({scrollTop: 0},400);  
            }  
              
            //返回底部  
            function returnBottom(){  
                $body.animate({scrollTop: $(document).height()},400);  
            }  
        </script>  
        <style type='text/css'>  
            #retopbtn{  
                position:fixed;  
                bottom:10px;  
                right:10px;  
            }  
            #rebtmbtn{  
                position:fixed;  
                top:10px;  
                right:10px;  
            }  
        </style>  
    </head>  
    <body>  
        <button id='rebtmbtn' onclick='returnBottom()'>⬇</button>  
		<button id='retopbtn' onclick='returnTop()'>⬆</button> 
	<body>


	    
	<div class="row my-card justify-content-center">
           
	<div class="col-lg-0 card">

	<!-- 上下翻转文字 -->
      
	<style type="text/css">#container-box-1{color:#526372;text-transform:uppercase;width:100%;font-size:16px;line-height:50px;text-align:center}#flip-box-1{overflow:hidden;height:50px}#flip-box-1 div{height:50px}#flip-box-1>div>div{color:#fff;display:inline-block;text-align:center;height:50px;width:100%}#flip-box-1 div:first-child{animation:show 20s linear infinite}.flip-box-1-1{background-color:#FF7E40}.flip-box-1-2{background-color:#C166FF}.flip-box-1-3{background-color:#737373}.flip-box-1-4{background-color:#4ec7f3}.flip-box-1-5{background-color:#42c58a}.flip-box-1-6{background-color:#F1617D}@keyframes show{0%{margin-top:-300px}5%{margin-top:-250px}16.666%{margin-top:-250px}21.666%{margin-top:-200px}33.332%{margin-top:-200px}38.332%{margin-top:-150px}49.998%{margin-top:-150px}54.998%{margin-top:-100px}66.664%{margin-top:-100px}71.664%{margin-top:-50px}83.33%{margin-top:-50px}88.33%{margin-top:0px}99.996%{margin-top:0px}100%{margin-top:300px}}</style>
	<div class="card card-site-info ">
	<div id="container-box-1">
	<div id="flip-box-1">
	<div><div class="flip-box-1-1"><i class="fa fa-gitlab" aria-hidden="true"></i>  rss feed for you </div></div>
	<div><div class="flip-box-1-2"><i class="fa fa-heart" aria-hidden="true"></i>  我们很年轻，但我们有信念、有梦想</div></div>
	<div><div class="flip-box-1-3"><i class="fa fa-gratipay" aria-hidden="true"></i>支持你的总会支持你，不支持的做再多也徒劳</div></div>
	<div><div class="flip-box-1-4"><i class="fa fa-drupal" aria-hidden="true"></i>  做这个世界的逆行者，先人一步看未来</div></div>
	<div><div class="flip-box-1-5"><i class="fa fa-gitlab" aria-hidden="true"></i>  只要你用心留意，世界将无比精彩</div></div>
	<div><div class="flip-box-1-6"><i class="fa fa-moon-o" aria-hidden="true"></i>  以下是信息聚合，精选各大站内容</div></div>
	<div><div class="flip-box-1-1">感谢原创者，感谢分享者，感谢值得尊重的每一位</div></div>
	</div>
	</div>
	</div>

			   <center>信息聚合阅读-RSS feed</center>
		
		<!-- 滚动代码-->

		<div class="card card-site-info ">
		<div class="m-3">
		<marquee scrollamount="5" behavior="right">
   
		<div id="blink">
   
		<a href="https://morss.it/:proxy:items=%7C%7C*[class=card]%7C%7Col%7Cli/https://rcy1314.github.io/news/" target="_blank">📢：rss feed for you 🔛</a>Rss聚合阅读页 🎁</div> 
   
   
		<script language="javascript"> 
   
   function changeColor(){ 
   
   var color="#f00|#0f0|#00f|#880|#808|#088|yellow|green|blue|gray"; 
   
   color=color.split("|"); 
   
   document.getElementById("blink").style.color=color[parseInt(Math.random() * color.length)]; 
   
   } 
   
   setInterval("changeColor()",200); 
   
		</script>
   
		</marquee>
		</div>
		</div>
   
   
		<!-- 向右流动代码-->
   
		<marquee scrollamount="3" direction="right" behavior="alternate">
   
		<a>😄😃😀</a>
   
		</marquee>
   
   
		
   
   
		<div class="alert alert-danger alert-dismissable">
		<button type="button" class="close" data-dismiss="alert"
			   aria-hidden="true">
		   &times;
		</button>
		 页面自动2小时监测更新一次！
		</div>
   
	<!-- 音乐 -->
	</script> 		  
	<div id="aplayer" class="aplayer" data-order="random" data-id="128460001" data-server="netease" data-type="playlist" data-fixed="true" data-autoplay="false" data-volume="0.8"></div>
	<!-- aplayer -->
	<script src="https://cdn.staticfile.org/jquery/3.2.1/jquery.min.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/aplayer@1.10.1/dist/APlayer.min.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/meting@1.2.0/dist/Meting.min.js"></script>
	<!-- end_aplayer -->
	<script src="https://cdn.staticfile.org/popper.js/1.15.0/umd/popper.min.js"></script>
	<script defer src="https://cdn.staticfile.org/twitter-bootstrap/4.4.1/js/bootstrap.min.js"></script>
	<script src="https://cdn.jsdelivr.net/gh/kaygb/kaygb@master/layer/layer.js"></script>
	<script src="https://cdn.jsdelivr.net/gh/kaygb/kaygb@master/js/v3.js"></script>
   
		<!-- 站长说 -->
   
		<div class="card card-site-info ">
		<div class="m-3">
		   <div class="small line-height-2"><b>公告 ： <i class="fa fa-volume-down fa-2" aria-hidden="true"></i></b></li><?php /*echo $conf['announcement'];*/?>  你可以点击上方rss feed for you来订阅页面，如需添加其它feed请点击页面最下方。</div>
		</div>
		 </div>
   
   
		<!-- 广告招租-->
		<div class="card card-site-info ">
		<div class="m-3">
		   <div class="small line-height-2"><b>广告位 <i class="fa fa-volume-down fa-2" aria-hidden="true"></i></b></li>：<?php /*echo $conf['announcement'];*/?>
		<a href="https://efficiencyfollow.notion.site" target="_blank" rel="nofollow noopener">
		<span>Efficiency主页</span></a>&nbsp;&nbsp;&nbsp; 
		<a href="https://noisedh.cn" target="_blank" rel="nofollow noopener">
		<span>Noise导航站</span></a>&nbsp;&nbsp;&nbsp;
		<a href="https://t.me/quanshoulu" target="_blank" rel="nofollow noopener">
		<span>TG发布频道</span></a>&nbsp;&nbsp;&nbsp;
		<a href="https://noisework.cn" target="_blank" rel="nofollow noopener">
		<span>引导主页</span></a>&nbsp;&nbsp;&nbsp;
		<a href="https://www.noisesite.cn" target="_blank" rel="nofollow noopener">
		<span>知识效率集</span></a>&nbsp;&nbsp;&nbsp;
		<a href="https://rcy1314.github.io/some-stars" target="_blank" rel="nofollow noopener">
		<span>我的star列表</span></a>&nbsp;&nbsp;&nbsp;
		<a href="https://noiseyp.top" target="_blank" rel="nofollow noopener">
		<span>Noise资源库</span></a></div>
		</div>
			<br>
	   

		<ol>
			{{ range .Posts }}<li><a href="{{ .Link }}" target="_blank" rel="noopener">{{ .Title }}</a> ({{ .Host }})</li>
			{{ end }}
		</ol>

		<footer>
		<div class="text-center py-1">   
        <div>
         <div class="text-center py-1">   
         <div>
		 <a href="https://ppnoise.notion.site/wiki-1ba2367142dc4b80b24873120a96efb5" target="_blank" rel="nofollow noopener">
	     <span>feed添加</span></a>    <br>
         </div>
	     <a href="https://noisework.cn" target="_blank" rel="nofollow noopener">
	     <span>主页</span></a>    <br>
         </div>
		 <script async src="//busuanzi.ibruce.info/busuanzi/2.3/busuanzi.pure.mini.js"></script>
		 <span id="busuanzi_container_site_pv" style='display:none'> 本站总访问量<span id="busuanzi_value_site_pv"></span>次</span>
		 </div>	
		 <div style="margin-top: 10px;">
		 &nbsp; 
		<span id="momk"></span>
		<span id="momk" style="color: #ff0000;"></span> 
		<script type="text/javascript">
   function NewDate(str) {
   str = str.split('-');
   var date = new Date();
   date.setUTCFullYear(str[0], str[1] - 1, str[2]);
   date.setUTCHours(0, 0, 0, 0);
   return date;
   }
   function momxc() {
   var birthDay =NewDate("2021-09-23");
   var today=new Date();
   var timeold=today.getTime()-birthDay.getTime();
   var sectimeold=timeold/1000
   var secondsold=Math.floor(sectimeold);
   var msPerDay=24*60*60*1000; var e_daysold=timeold/msPerDay;
   var daysold=Math.floor(e_daysold);
   var e_hrsold=(daysold-e_daysold)*-24;
   var hrsold=Math.floor(e_hrsold);
   var e_minsold=(hrsold-e_hrsold)*-60;
   var minsold=Math.floor((hrsold-e_hrsold)*-60); var seconds=Math.floor((minsold-e_minsold)*-60).toString();
   document.getElementById("momk").innerHTML = "本站已运行:"+daysold+"天"+hrsold+"小时"+minsold+"分"+seconds+"秒";
   setTimeout(momxc, 1000);
   }momxc();
	</footer>
</body>
</html>
`

	tmpl, err := template.New("webpage").Parse(htmlTemplate)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(writer, templateData); err != nil {
		return err
	}

	return nil
}