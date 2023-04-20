rule GC2Detector
{

	strings:
		$gc2 = "gc2-sheet"

	condition:
		$gc2

}

rule GC2Authentication {

	strings:
		$authenticationSheet = "AuthenticateSheet"
		$authenticationDrive = "AuthenticateDrive"

	condition:
		any of them

}

rule GC2Create {

	strings:
		$Create = "createSheet"

	condition:
		$Create

}

rule GC2Download {

	strings:
		$Download = "downloadFile"

	condition:
		$Download

}

rule GC2Read {

	strings:
		$Read = "readSheet"

	condition:
		$Read
}

rule GC2Write {

	strings:
		$Write = "writeSheet"

	condition:
		$Write

}

rule GC2Options {

	strings:
		$GetOptionsCredential = "GetOptionsCredential"
		$GetOptionsSheetId = "GetOptionsSheetId"
		$GetOptionsDriveId = "GetOptionsDriveId"
		$GetOptionsDebug = "GetOptionsDebug"

	condition:
		any of them

}

rule GC2Utils {

	strings:
		$GenerateNewSheetName = "GenerateNewSheetName"
		$GetLastCommand = "GetLastCommand"
		$CreateNewEmptyCommand = "CreateNewEmptyCommand"

	condition:
		any of them

}

rule GC2Struct {
	strings:
		$RangeTickerConfiguration = "RangeTickerConfiguration"

	condition:
		$RangeTickerConfiguration
}

rule GC2Confirmation {
	condition:
		GC2Authentication and GC2Struct and GC2Utils and GC2Options
}