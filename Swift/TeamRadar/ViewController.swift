/*
Copyright 2016 IslandJohn and the TeamRadar Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied
See the License for the specific language governing permissions and
limitations under the License.
*/

import Cocoa

class ViewController: NSViewController {

    @IBOutlet weak var prefUrlText: NSTextField!
    @IBOutlet weak var prefUserText: NSTextField!
    @IBOutlet weak var prefPasswordText: NSSecureTextField!
    @IBOutlet weak var prefSaveButton: NSButton!
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }

    override var representedObject: AnyObject? {
        didSet {
        // Update the view, if already loaded.
        }
    }
    
    @IBAction func cancelAction(sender: AnyObject) {
    }
    
    @IBAction func saveAction(sender: AnyObject) {
    }
}
