#!/usr/bin/env python3
"""Creates the Automator Quick Action for Clipboard History."""
import plistlib, pathlib

service_dir = pathlib.Path.home() / "Library/Services/Clipboard History.workflow/Contents"
service_dir.mkdir(parents=True, exist_ok=True, mode=0o700)

workflow = {
    "AMApplicationBuild": "512",
    "AMApplicationVersion": "2.10",
    "AMDocumentVersion": "2",
    "actions": [
        {
            "action": {
                "AMAccepts": {"Container": "List", "Optional": True, "Types": ["com.apple.cocoa.string"]},
                "AMActionVersion": "2.0.3",
                "AMApplication": ["Any Application"],
                "AMParameterProperties": {
                    "COMMAND_STRING": {}, "CheckedForUserDefaultShell": {},
                    "inputMethod": {}, "shell": {}, "source": {}
                },
                "AMProvides": {"Container": "List", "Types": ["com.apple.cocoa.string"]},
                "ActionBundlePath": "/System/Library/Automator/Run Shell Script.action",
                "ActionName": "Run Shell Script",
                "ActionParameters": {
                    "COMMAND_STRING": "/usr/bin/python3 ~/.clipboard-history/chooser.py",
                    "CheckedForUserDefaultShell": True,
                    "inputMethod": 0,
                    "shell": "/bin/bash",
                    "source": ""
                },
                "BundleIdentifier": "com.apple.RunShellScript",
                "CFBundleVersion": "2.0.3",
                "CanShowSelectedItemsWhenRun": False,
                "CanShowWhenRun": True,
                "Category": ["AMCategoryUtilities"],
                "Class Name": "RunShellScriptAction",
                "InputUUID": "auto-input",
                "Keywords": ["Shell", "Script", "Command", "Run", "Unix"],
                "OutputUUID": "auto-output",
                "UUID": "auto-action",
                "UnlocalizedApplications": ["Automator"],
                "arguments": {},
                "isViewVisible": True,
                "location": "530.5:242",
                "nickname": "Run Shell Script",
                "overrideTextColor": False,
                "runtimeProperties": {"RequiresInput": False}
            }
        }
    ],
    "connectors": {},
    "workflowMetaData": {
        "applicationBundleIDsByPath": {},
        "applicationPaths": [],
        "inputTypeIdentifier": "com.apple.Automator.nothing",
        "outputTypeIdentifier": "com.apple.Automator.nothing",
        "presentationMode": 11,
        "processesInput": False,
        "serviceInputTypeIdentifier": "com.apple.Automator.nothing",
        "serviceOutputTypeIdentifier": "com.apple.Automator.nothing",
        "serviceProcessesInput": False,
        "systemImageName": "NSActionTemplate",
        "useAutomaticInputType": False,
        "workflowTypeIdentifier": "com.apple.Automator.servicesMenu"
    }
}

with open(service_dir / "document.wflow", "wb") as f:
    plistlib.dump(workflow, f)

info = {
    "CFBundleIdentifier": "com.apple.automator.Clipboard-History",
    "CFBundleName": "Clipboard History",
    "CFBundlePackageType": "APPL",
    "CFBundleShortVersionString": "1.0",
    "CFBundleVersion": "1",
    "NSPrincipalClass": "AMWorkflowController",
    "NSServices": [
        {
            "NSMenuItem": {"default": "Clipboard History"},
            "NSMessage": "runWorkflowAsService",
            "NSRequiredContext": {},
            "NSSendTypes": []
        }
    ]
}
with open(service_dir / "Info.plist", "wb") as f:
    plistlib.dump(info, f)

print("Done:", service_dir.parent)
