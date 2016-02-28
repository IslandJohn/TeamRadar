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

    @IBOutlet var radarMenu: NSMenu!
    @IBOutlet var radarStateMenuItem: NSMenuItem!
    
    var radarStatusItem: NSStatusItem? = nil
    var radarGoTask: NSTask? = nil
    
    override func viewDidLoad() {
        super.viewDidLoad()

        radarStatusItem = NSStatusBar.systemStatusBar().statusItemWithLength(NSVariableStatusItemLength)
        
        radarStatusItem?.title = "TeamRadar"
        radarStatusItem?.highlightMode = true
        radarStatusItem?.menu = radarMenu
    }

    override var representedObject: AnyObject? {
        didSet {
        // Update the view, if already loaded.
        }
    }
    
    func eventTaskOutput(note: NSNotification) {
        let fh = note.object as! NSFileHandle
        
        fh.waitForDataInBackgroundAndNotify()
    }
    
    func eventTaskError(note: NSNotification) {
        let fh = note.object as! NSFileHandle
        
        fh.waitForDataInBackgroundAndNotify()
    }
    
    func eventTaskTerminate() {
    }
    
    @IBAction func connectAction(sender: AnyObject) {
        let menuitem = sender as? NSMenuItem
        
        if (radarGoTask == nil || !radarGoTask!.running) {
            if (radarGoTask == nil) {
                radarGoTask = NSTask()
                
                radarGoTask!.launchPath = NSBundle.mainBundle().pathForResource("teamradar", ofType: nil)
                radarGoTask!.standardInput = NSPipe()
                radarGoTask!.standardOutput = NSPipe()
                radarGoTask!.standardError = NSPipe()
                radarGoTask?.terminationHandler = {(task: NSTask) -> Void in
                    self.eventTaskTerminate()
                }
                
                NSNotificationCenter.defaultCenter().addObserver(self, selector: "eventTaskOutput:", name: NSFileHandleDataAvailableNotification, object: radarGoTask!.standardOutput?.fileHandleForReading)
                NSNotificationCenter.defaultCenter().addObserver(self, selector: "eventTaskError:", name: NSFileHandleDataAvailableNotification, object: radarGoTask!.standardError?.fileHandleForReading)
                
                radarGoTask!.standardOutput?.fileHandleForReading.waitForDataInBackgroundAndNotify()
                radarGoTask!.standardError?.fileHandleForReading.waitForDataInBackgroundAndNotify()
            }
            
            //radarGoTask?.arguments = nil
            radarGoTask?.launch()
            
            menuitem?.title = "Disconnect"
            radarStateMenuItem.title = "No rooms."
        }
        else {
            radarGoTask?.waitUntilExit()
            
            menuitem?.title = "Connect..."
            radarStateMenuItem.title = "Not connected."
        }
    }
}

