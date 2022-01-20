import { Application } from "@hotwired/stimulus"
import ClipboardController from './controllers/clipboard_controller'
const application = Application.start()
application.register('clipboard', ClipboardController)
