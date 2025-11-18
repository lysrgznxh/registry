package thread

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/util/log"
	"fmt"
)

// 安全运行异步协程，当发生错误时，回调err方法
func NewThread(name string, runCtrl *core.RuntimeControl, fn func(runCtrl *core.RuntimeControl), err func(err string)) {
	log.Info("Starting ", name)
	runCtrl.Wg.Add(1)
	go runSafe(name, runCtrl, fn, err)
}

// 安全运行异步协程，当发生错误时，回调err方法
func runSafe(name string, runCtrl *core.RuntimeControl, fn func(runCtrl *core.RuntimeControl), errCall func(err string)) {
	defer func() {
		runCtrl.Wg.Done() //运行结束时给主进程信号 -1 ,为什么要放在这里? 因为只有放在这里才能确保一定会释放计数器
		recover_object := recover()
		if runCtrl.IsCancel { //如果是主线程释放的退出,则这里不回调错误
			return
		}
		err, ok := recover_object.(error)
		if ok { //如果是错误类型
			errCall(err.Error())
			//err = errors.New("")
		} else { //不是错误类型,
			errCall(fmt.Sprintf("%v", recover_object)) //加一句是为了，当线程结束时也能通知给主线程
		}
	}()
	fn(runCtrl)
	//以下代码只有正常结束才会执行,如果fn执行的时候发生了pansic则只会执行defer func内部的代码
	log.Error(name, " exited!")
}

func NewThreadV2(name string, runCtrl *core.RuntimeControl, fn func(runCtrl *core.RuntimeControl)) {
	runCtrl.ThreadSize += 1
	runCtrl.AddThreadLog(name, 1)
	runCtrl.Wg.Add(1)
	go runSafeV2(name, runCtrl, fn)
}

// 安全运行异步协程，当发生错误时，关闭主线程
func runSafeV2(name string, runCtrl *core.RuntimeControl, fn func(runCtrl *core.RuntimeControl)) {
	defer func() {
		runCtrl.Wg.Done() //运行结束时给主进程信号 -1 ,为什么要放在这里? 因为只有放在这里才能确保一定会释放计数器
		recover_object := recover()
		if runCtrl.IsCancel { //如果是主线程释放的退出,则这里不回调错误
			return
		}
		var errContent = ""
		err, ok := recover_object.(error)
		if ok { //如果是错误类型
			errContent = name + " Exited! error content:" + err.Error()
		} else { //不是错误类型,
			errContent = fmt.Sprintf("%s Exited! error content::%v", name, recover_object)
		}
		log.Error(errContent)
		runCtrl.Close(errContent) //关闭主线程
	}()
	fn(runCtrl)
	//以下代码只有正常结束才会执行,如果fn执行的时候发生了pansic则只会执行defer func内部的代码
	runCtrl.AddThreadLog(name, 0)
}
