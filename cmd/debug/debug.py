# after Safari 16.4!
import os
import frida


def check() -> frida.core.Device:
    dev = frida.get_local_device()
    if dev.query_system_parameters()['os']['id'] != 'macos':
        raise RuntimeError('This script is only for Mac OS')

    if "enabled" in os.popen("csrutil status").read():
        raise RuntimeError('SIP must be closed')
    return dev


def hook(dev: frida.core.Device):
    for proc in dev.enumerate_processes():
        if proc.name not in [
            'WeChat', 'Mini Program',  # en
            '微信', '小程序'  # zh-cn
        ]:
            continue

        print('Patching %s (%d)' % (proc.name, proc.pid))
        session = dev.attach(proc.pid)
        script = session.create_script('''
            ['WKWebView', 'JSContext'].forEach(
                clazz => ObjC.chooseSync(ObjC.classes[clazz]).forEach(
                    v => v.setInspectable_(ptr(1))
                )
            )
        ''')
        script.load()
        script.unload()
        session.detach()


if __name__ == '__main__':
    dev = check()
    hook(dev)
