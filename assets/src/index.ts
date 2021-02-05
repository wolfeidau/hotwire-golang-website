import { Application } from "@stimulus/core"
import ClipboardController from './controllers/clipboard_controller'
const application = Application.start()
application.register('clipboard', ClipboardController)
