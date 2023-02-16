package device

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDevice_GetBlackbox(t *testing.T) {
	d, err := NewBuilder("test1").
		SetOsName(MacOSX).
		SetBrowserName(Chrome).
		SetMemory(8).
		SetHardwareConcurrency(16).
		ScreenColorDepth(24).
		SetScreenWidth(1900).
		SetScreenHeight(900).
		SetTimezone("America/Los_Angeles").
		SetLanguages("en-US,en").
		Build()
	if err != nil {
		panic(err)
	}
	fprnt, _ := d.GetBlackbox()
	by, _ := json.Marshal(fprnt)
	fmt.Println(string(by))
	assert.Equal(t, 1, 2)
}

func Test_DecryptBlackbox(t *testing.T) {
	expected := "[8," +
		"\"America/Los_Angeles\"," +
		"false," +
		"\"Blink\"," +
		"\"Mac OS X\"," +
		"\"Chrome\"," +
		"\"Google Inc.\"," +
		"8," +
		"12," +
		"\"en-US,fr-CA\"," +
		"\"f473d473013d58cee78732e974dd4af2e8d0105449c384658cbf1505e40ede50\"," +
		"\"Google Inc. (ATI Technologies Inc.),ANGLE (ATI Technologies Inc., AMD Radeon Pro 5300M OpenGL Engine, OpenGL 4.1)\"," +
		"\"e069294a1b55cc6f1dfb0ae60a7f104a9893cbaa8d94772f7f11f48b04448eea\"," +
		"\"d9af7aa1d00f202e8291fe49b9344f69746635eea53e7eace68c10f302cc933a\"," +
		"1920," +
		"991," +
		"24," +
		"true," +
		"true," +
		"\"1f03b77fda33742261bea0d27e6423bf22d2bf57febc53ae75b962f6e523cc02\"," +
		"\"5a0ef26fd9ff096689feaad0d49fb8551822ea6b3be74a02794c2aa10ead141f\"," +
		"\"6aeb6412b24ba7dd08653eb50179026602499917a6400174f9ad7e9bef78abf2\"," +
		"\"5e78318249d4fd2930a116b628631f99ef9068c8f2e34940330389796bd9a0ba\"," +
		"124.04347657808103," +
		"\"d563c9730f2852b84227da38496ac56130adc9a284b4c10b673fd2a43781ee70\"," +
		"1561173311," +
		"\"2023-02-09T08:42:07.132Z\"," +
		"\"l0469ofavup5il806wnm01ni9mz\"," +
		"272," +
		"\"10_15_7\"," +
		"\"RnNEfC8iU0tNQ1x2MVtVXkRjSyZfViNiYitjZjVuOjEjXkFhTUVVcCxJVSVaQ3BjfFRpWXZ1cWk2ZHpPOnlDRUUmYUdRQFNmIWpRU3U5ei9LcWFQOkNCQTZKcWA4JlFdM1JaKSAxNjc1OTMyMTI2NDkz\"," +
		"\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36\"," +
		"\"2023-02-09T08:42:07.000Z\"," +
		"null]"
	decrypted, _ := DecryptBlackbox("JVqc1PkrbpPF9zilCnzlSKnOAEaSAXTTFILpTrofkrfpG0BytRt86FvA5Rdaf7HjJZH6aNP4KlyBs_YbTX_MLZC15xdmud4QQJi97yFGeLvgEkSH72HQPaLH-StQgsXqHE6VBHPaRqvQAjJ76Ux6n9EDKFqd1fosb6DS9ylskcP1Wsj1Sp3C9DedDzx_wOUXSW6g4wg6bNIGPXDUCD9yotMGap_XOp8EO3Oq3Q90reQYfOAUddsNcqoOPm-f1Ag8ddgLQ3et4hp930V2q9sQdanZPqIHPGyRw_UaTI-05hhfzj2kEHWazPxFsxZEaZvL8zSI0fYoWKwRdNxKuSWU-2TJPGGTwwx63Qs0WYvOD12k8DVajLzkJXnC5xlJnQJlzTuqFoXsVbotUoS0_WvO_CFTlrvtHV6r7xRGdsgpjfJhz_QmVqYYh6zeDkN2ptYjSHqq-WnOPIPP9CZWmwlw2Ues0QNGa53NHIzxX6byF0l5rdsMNVqMvuMVWH2v4UZ2rOUXUITlFnit4kWo3kR12T-h0TKXzf1elfssXJDxKmKbzjGT9FWN8Spelcz-ZJsBMmPJ_TWXx_svY5sAZcbrHU90pukOQHLWD3DWDW7PAGSUxCpcjL4jW43G913C9i-Ryv0xZcsBOnGl2xFEed5DpNkMcagNbtE2bKQHOGjOATFjxililcgpToCy1wlMfbboGD1vsuskVXqs7yFVeqzvY9VKr9QGSb0vpAkuYKPI-ixdw_MmiL_2XMAhVIe-8iRWjL0fhOUVeaviR32x4xZ43hBCptg6oNUMctc5nNEEZcoBNpjRBzmf1TpvodQ3msr8IVOFqtwfRHao3T5u0zlroQdrpApwoNkPRX22HIHiQ6fXO2-oDnCo3RJDe63fRKXbPXDSN26iAzNlnNUJbJ7_YJHBJofrHFCB5ww-cJXHCi9hk8kqj_EnW4y-IFKG6EmA5Eh4sOYbTrMVSnqr4htLfbPpGUt_uPEqW5LzKV2Nve4lWb_4Wb30WZL0Wb_2Lo_xV4mu4BI3aazRAzVqzwY-caLaDEB53RF32w1GeakKO2yiBDpspNoNPqTdFnvhGkqAuBtTuetQg7fwJFSHuuodVY7F_jSW-jOUxCaHrN4QNWeq2w1Bb5_TBjpxp9wTS3uz5BRHbJ7hBjhqzgM5bM8IP3KiCDpyp9k7c6fZC0KmBzpypt8VdtkORHWo2DmdADmazAQ4ms4xYpL0KmGU-l6Q8SVYj8f4XcL5KU6AstcJTH2y6BlKgbTnGEluoOMIOmyezgAzYJDC7x9YrNwUOWyt4RM4a6zcE0FypdcxVoi63xFUeavdSXmt4xyL8VLIPa3iS7fvH1XMOqfXCHbfGIX_JFaIrd8iVIu94hRXfK7gEUGg0QZlnMHzJUp8v-QWSJoIVpsBRHzlOmreLH2uJlil-2_FHYjaRJcQatAmj91Gnwh85kCqAHXELnPdNaDmTqL3TaMGScELYbQKa7zvMZsBR5kJYLgSQ6b9aJr0PKz8S7klabsQZdIrgOQ2h80biNEomOo_csf8YcoDT7IJT6DvWqjrPJDqNZjvMGSuGmDEEUKM7TiLzESS_F-Q3zOA-Uaa4xVjpxKMseMVOmyv1AY4hfRu10OvEDVnreIQQGWXx-88nQBp10u6LZW67S9Uhrb_beFGstcJOYbnSm-h0SBzmMr6Unep2Qo6mcr_XpW-4xVFhvZm0jeO81WgCX2i1BpPgrnnGlB1p9f_SpLmM3-k1hk-cKAMdeBFapzME3jbRrXeAzVlqBCC8V7D6BpgkcH6KFiGtuQUOWub7k-1FojxFkiOw_YtW47E6RtNcqTnDD5wotIEN2SUxvMjXLDgGD1wseUXPG-w4BdFdaXVL1SGuN0PUsA1oQ0yZ6s")
	decrypted2, _ := DecryptBlackbox(EncryptBlackbox(decrypted))
	assert.Equal(t, expected, decrypted)
	assert.Equal(t, expected, decrypted2)

	assert.Equal(t, "dA", EncryptBlackbox("t"))
	assert.Equal(t, "dNk", EncryptBlackbox("te"))
	assert.Equal(t, "dNlM", EncryptBlackbox("tes"))
	assert.Equal(t, "dNlMwA", EncryptBlackbox("test"))
	assert.Equal(t, "dNlMwPE", EncryptBlackbox("test1"))
}
