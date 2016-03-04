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

class ViewController: NSViewController, NSUserNotificationCenterDelegate {

    @IBOutlet weak var prefUrlText: NSTextField!
    @IBOutlet weak var prefUserText: NSTextField!
    @IBOutlet weak var prefPasswordText: NSSecureTextField!
    @IBOutlet weak var prefSaveButton: NSButton!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        prefUrlText.stringValue = Settings.get(SettingsKey.SERVER) ?? ""
        prefUserText.stringValue = Settings.get(SettingsKey.USERNAME) ?? ""
        prefPasswordText.stringValue = Settings.get(SettingsKey.PASSWORD) ?? ""
    }

    override var representedObject: AnyObject? {
        didSet {
        // Update the view, if already loaded.
        }
    }
    
    @IBAction func cancelAction(sender: AnyObject) {
        [self .dismissViewController(self)];
    }
    
    @IBAction func saveAction(sender: AnyObject) {
        Settings.set(SettingsKey.USERNAME, value: prefUserText.stringValue)
        Settings.set(SettingsKey.PASSWORD, value: prefPasswordText.stringValue)
        Settings.set(SettingsKey.SERVER, value: prefUrlText.stringValue)
    }
    
    func showNotification() -> Void {
        let unc = NSUserNotificationCenter.defaultUserNotificationCenter()
        unc.delegate = self
        let notification = NSUserNotification()
        notification.title = "Test from Swift"
        notification.informativeText = "The body of this Swift notification"
        unc.deliverNotification(notification)
    }
    
    func userNotificationCenter(center: NSUserNotificationCenter, shouldPresentNotification notification: NSUserNotification) -> Bool {
        return true
    }
}

