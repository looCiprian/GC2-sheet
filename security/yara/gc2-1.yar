rule GC2Detector
{

	strings:
		$authenticationSheet = "AuthenticateSheet"
		$authenticationDrive = "AuthenticateDrive"
		$Read = "readSheet"
		$Write = "writeSheet"
		$GetOptionsCredential = "GetOptionsCredential"
		$GetOptionsSheetId = "GetOptionsSheetId"
		$GetOptionsDriveId = "GetOptionsDriveId"
		$GenerateNewSheetName = "GenerateNewSheetName"
		$RangeTickerConfiguration = "RangeTickerConfiguration"

	condition:
		(any of ($RangeTickerConfiguration, $GetOptionsDriveId, $GetOptionsSheetId)) or ($Read and $Write) or $RangeTickerConfiguration or (any of ($authenticationSheet, $authenticationDrive)) or $GetOptionsCredential or $GenerateNewSheetName
}