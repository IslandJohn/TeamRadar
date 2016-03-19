//
//  TeamRadarTests.swift
//  TeamRadarTests
//
//  Created by Janos on 2/27/16.
//  Copyright © 2016 IslandJohn. All rights reserved.
//

import XCTest
@testable import TeamRadar

class TeamRadarTests: XCTestCase {
    
    override func setUp() {
        super.setUp()
        // Put setup code here. This method is called before the invocation of each test method in the class.
    }
    
    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
        super.tearDown()
    }
    
    func testExtractJsonFromLine() {
        let str = "messages new 8214 828bab5f-0f36-4e8f-bbc2-e98cbfe67e8f 6336 {\"Id\":6336,\"Content\":\"asdjfsj asdfjasd asjdf\",\"MessageType\":\"normal\",\"PostedTime\":\"2016-03-16T00:42:29.367Z\",\"PostedRoomId\":8214,\"PostedBy\":{\"Id\":\"828bab5f-0f36-4e8f-bbc2-e98cbfe67e8f\",\"DisplayName\":\"Felix\",\"Url\":\"https://islandjohn.vssps.visualstudio.com/_apis/Identities/828bab5f-0f36-4e8f-bbc2-e98cbfe67e8f\",\"ImageUrl\":\"https://islandjohn.visualstudio.com/DefaultCollection/_api/_common/identityImage?id=828bab5f-0f36-4e8f-bbc2-e98cbfe67e8f\"}}"
        let justJson = "{\"Id\":6336,\"Content\":\"asdjfsj asdfjasd asjdf\",\"MessageType\":\"normal\",\"PostedTime\":\"2016-03-16T00:42:29.367Z\",\"PostedRoomId\":8214,\"PostedBy\":{\"Id\":\"828bab5f-0f36-4e8f-bbc2-e98cbfe67e8f\",\"DisplayName\":\"Felix\",\"Url\":\"https://islandjohn.vssps.visualstudio.com/_apis/Identities/828bab5f-0f36-4e8f-bbc2-e98cbfe67e8f\",\"ImageUrl\":\"https://islandjohn.visualstudio.com/DefaultCollection/_api/_common/identityImage?id=828bab5f-0f36-4e8f-bbc2-e98cbfe67e8f\"}}"
        
        let parser = TeamRadarParser()
        let json = parser.extractJSONFromLine(str)
        
        XCTAssertEqual(json, justJson)
        
        let jsonDict: AnyObject? = parser.convertJSONStringToDictionary(json)
        
        guard let jDict = jsonDict else { XCTFail("JSON Dictionary was nil"); return }
        guard jDict is NSDictionary else { XCTFail("JSON is not a NSDictionary"); return }
        
        let jsonDictionary = jDict as! NSDictionary
        
        XCTAssertEqual(jsonDictionary["MessageType"] as? String, "normal")
        
        let postedByDict: NSDictionary = jsonDictionary["PostedBy"] as! NSDictionary
        
        XCTAssertEqual(postedByDict["DisplayName"] as? String, "Felix")
        
    }
    
}
