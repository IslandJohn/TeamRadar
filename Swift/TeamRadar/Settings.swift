//
//  File.swift
//  TeamRadar
//
//  Created by Felix Khazin on 3/3/16.
//  Copyright © 2016 IslandJohn. All rights reserved.
//

import Foundation

enum SettingsKey: String {
    case USERNAME = "username"
    case PASSWORD = "password"
    case SERVER = "server"
}

class Settings {
    static func get(key: SettingsKey) -> String?{
        return NSUserDefaults.standardUserDefaults().stringForKey(key.rawValue)
    }
    
    static func set(key:SettingsKey, value:String){
        NSUserDefaults.standardUserDefaults().setObject(value, forKey: key.rawValue)
        NSUserDefaults.standardUserDefaults().synchronize()
    }
}