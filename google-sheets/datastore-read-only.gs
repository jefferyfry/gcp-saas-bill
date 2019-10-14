var CONFIG = "datastore-viewer-service-account.json";

var datastore = {
    /* api related */
    scopes:        "https://www.googleapis.com/auth/datastore https://www.googleapis.com/auth/drive",
    baseUrl:       "https://datastore.googleapis.com/v1",
    httpMethod:    "POST",
    currentUrl:    false,

    /* authentication */
    oauth:         false,
    projectId:     false,
    clientId:      false,
    clientEmail:   false,
    privateKey:    false,

    getInstance: function(){
        if(! this.config) {this.getConfig(CONFIG);}

        if(! this.oauth) {this.getService();}

        return this;
    },

    getConfig: function(filename){
        var it = DriveApp.getFilesByName(filename);
        while (it.hasNext()) {
            var file = it.next();
            var data = JSON.parse(file.getAs("application/json").getDataAsString());
            this.projectId   = data.project_id;
            this.privateKey  = data.private_key;
            this.clientEmail = data.client_email;
            this.clientId    = data.client_id;
            continue;
        }
        this.currentUrl = this.baseUrl + "/projects/" + this.projectId + ":runQuery";
    },

    getService: function(){
        this.oauth = OAuth2.createService("Datastore")
            .setTokenUrl("https://www.googleapis.com/oauth2/v4/token")
            .setPropertyStore(PropertiesService.getScriptProperties())
            .setPrivateKey(this.privateKey)
            .setIssuer(this.clientEmail)
            .setScope(this.scopes);
    },

    authCallback: function(request) {
        var service = createService();
        var isAuthorized = service.handleCallback(request);
        if (isAuthorized) {
            return HtmlService.createHtmlOutput('Success! You can close this tab.');
        } else {
            return HtmlService.createHtmlOutput('Denied. You can close this tab');
        }
    },

    getOptions: function(payload) {
        return {
            headers: {Authorization: 'Bearer ' + this.oauth.getAccessToken()},
            payload: JSON.stringify(payload),
            contentType: "application/json",
            muteHttpExceptions: true,
            method: this.httpMethod
        };
    },

    runGql: function(query_string) {
        if(! this.startCursor) {
            var gqlOpt = {
                gqlQuery: {
                    query_string: query_string,
                    allowLiterals: true
                }
            };
        } else {
            var gqlOpt = {
                gqlQuery: {
                    query_string: query_string,
                    allowLiterals: true,
                    namedBindings: {
                        startCursor: {cursor: this.startCursor}
                    }
                }
            };
        }

        /* wait for the previous request to complete, check every 200 ms */
        while(this.queryInProgress) {Utilities.sleep(200);}

        /* set queryInProgress flag to true */
        this.queryInProgress = true;

        /* configure the request */
        var options = this.getOptions(gqlOpt);

        /* execute the request */
        var response = UrlFetchApp.fetch(this.currentUrl, options);
        var result = JSON.parse(response.getContentText());

        /* set queryInProgress flag to false */
        this.queryInProgress = false;

        /* always log remote errors */
        if(typeof(result.error) !== "undefined") {
            Logger.log(method + " > ERROR " + result.error.code + ": " + result.error.message);
            return false;
        }

        if(typeof(result.batch) !== "undefined")
            return result.batch['entityResults'];
        else
            return false;
    },
}



function main(){
    var currentdate = new Date();
    var datetime = "Last Sync: " + currentdate.getDate() + "/"
        + (currentdate.getMonth()+1)  + "/"
        + currentdate.getFullYear() + " @ "
        + currentdate.getHours() + ":"
        + currentdate.getMinutes() + ":"
        + currentdate.getSeconds();

    //clear contents
    SpreadsheetApp.getActiveSpreadsheet().getActiveSheet().clear();

    //set last sync date/time
    var dateTime = [[datetime]];
    SpreadsheetApp.getActiveSpreadsheet().getActiveSheet().getRange("A1:A1").setValues(dateTime);

    //set header
    var header = [["Account Id","Account Status","First Name","Last Name","Email","Phone","Timezone","Entitlement ID","Product","Plan","Entitlement Status"]];
    SpreadsheetApp.getActiveSpreadsheet().getActiveSheet().getRange("A2:K2").setValues(header);

    //query accounts
    var accounts = datastore.getInstance().runGql("select * from Account");
    if (accounts === undefined || accounts.length == 0)
        return;
    var row = 3;
    for(i=0; i < accounts.length; i++) {
        var account = accounts[i];
        var accountId = account['entity']['properties']['id']['stringValue'];
        var contacts = datastore.getInstance().runGql("select * from Contact where accountId ='"+accountId+"'");

        if (contacts === undefined || contacts.length == 0)
            continue;

        var contact = contacts[0];

        var entitlements = datastore.getInstance().runGql("select * from Entitlement where account ='"+accountId+"'");

        if (entitlements === undefined || entitlements.length == 0)
            continue;

        var accountStatus = account['entity']['properties']['state']['stringValue'];

        var firstName = contact['entity']['properties']['firstName']['stringValue'];
        var lastName = contact['entity']['properties']['lastName']['stringValue'];
        var emailAddress = contact['entity']['properties']['emailAddress']['stringValue'];
        var phone = contact['entity']['properties']['phone']['stringValue'];
        var timezone = contact['entity']['properties']['timezone']['stringValue'];
        for(j=0; j < entitlements.length; j++) {
            var entitlement = entitlements[j];
            var entitlementId = entitlement['entity']['properties']['id']['stringValue'];
            var product = entitlement['entity']['properties']['product']['stringValue'];
            var plan = entitlement['entity']['properties']['plan']['stringValue'];
            var state = entitlement['entity']['properties']['state']['stringValue'];
            var rowValues = [[accountId,accountStatus,firstName,lastName,emailAddress,phone,timezone,entitlementId,product,plan,state]];
            SpreadsheetApp.getActiveSpreadsheet().getActiveSheet().getRange("A"+row+":K"+row).setValues(rowValues);
            row++;
        }
        SpreadsheetApp.getActiveSpreadsheet().getActiveSheet().autoResizeColumns(1, 11);
    }
}