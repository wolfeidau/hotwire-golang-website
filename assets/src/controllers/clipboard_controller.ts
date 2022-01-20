import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  sourceTarget: HTMLInputElement | undefined
  static targets = ["source"]

  connect() {
    console.log('connect')
  }

  copy() {
    this.sourceTarget?.select()
    document.execCommand("copy")
  }
}
