['WKWebView', 'JSContext'].forEach(
    clazz => ObjC.chooseSync(ObjC.classes[clazz]).forEach(
        v => v.setInspectable_(ptr(1))
    )
)