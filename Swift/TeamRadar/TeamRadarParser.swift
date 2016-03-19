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

import Foundation

public class TeamRadarParser: NSObject {
    func extractJSONFromLine(line:String) -> String{
        do {
            let regex = try NSRegularExpression(pattern: "(\\{.*\\})", options: .CaseInsensitive)
            let result: NSTextCheckingResult? = regex.firstMatchInString(line, options: .ReportCompletion, range: NSMakeRange(0, line.characters.count))
            
            guard let r = result else { return "" }
            
            let rangeOfJson: NSRange? = r.rangeAtIndex(1)
            
            guard let rJson = rangeOfJson else { return "" }
            
            let nsLine = line as NSString
            
            let jsonString = nsLine.substringWithRange(rJson)
            
            return jsonString
        } catch {
            return ""
        }
    }
    
    func convertJSONStringToDictionary(json:String) -> AnyObject? {
        let encoding = NSUTF8StringEncoding
        let jsonData = json.dataUsingEncoding(encoding)
        guard let jData = jsonData else {return nil}
        do {
            return try NSJSONSerialization.JSONObjectWithData(jData, options: [])
        } catch let error {
            print("json error: \(error)")
            return nil
        }
    }
}
