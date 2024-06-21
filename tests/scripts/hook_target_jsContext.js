// https://webkit.org/blog/13936/enabling-the-inspection-of-web-content-in-apps/
['WKWebView', 'JSContext'].forEach(
    clazz => ObjC.chooseSync(ObjC.classes[clazz]).forEach(
        v => v.setInspectable_(ptr(1))
    )
)
