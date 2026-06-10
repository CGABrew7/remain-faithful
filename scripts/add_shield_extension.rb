require 'xcodeproj'

PROJECT_PATH = File.expand_path('../../RemainFaithful.xcodeproj', __FILE__)
project = Xcodeproj::Project.open(PROJECT_PATH)

TARGET_NAME   = 'RemainFaithfulShieldConfig'
BUNDLE_ID     = 'com.remainfaithful.app.ShieldConfig'
EXTENSION_DIR = File.expand_path('../../RemainFaithfulShieldConfig', __FILE__)

# ── Check if already added ──────────────────────────────────────────────────
if project.targets.any? { |t| t.name == TARGET_NAME }
  puts "Target #{TARGET_NAME} already exists — skipping"
  exit 0
end

main_target = project.targets.find { |t| t.name == 'RemainFaithful' }
abort('Cannot find RemainFaithful target') unless main_target

# ── Create extension target ─────────────────────────────────────────────────
ext_target = project.new_target(
  :app_extension,
  TARGET_NAME,
  :ios,
  '17.0'
)

# ── Build settings ───────────────────────────────────────────────────────────
['Debug', 'Release'].each do |config_name|
  config = ext_target.build_configuration_list.build_settings(config_name)
  config['PRODUCT_NAME']              = TARGET_NAME
  config['PRODUCT_BUNDLE_IDENTIFIER'] = BUNDLE_ID
  config['INFOPLIST_FILE']            = "#{TARGET_NAME}/Info.plist"
  config['CODE_SIGN_ENTITLEMENTS']    = "#{TARGET_NAME}/#{TARGET_NAME}.entitlements"
  config['SWIFT_VERSION']             = '5.0'
  config['IPHONEOS_DEPLOYMENT_TARGET'] = '17.0'
  config['TARGETED_DEVICE_FAMILY']    = '1,2'
  config['SKIP_INSTALL']              = 'YES'
  config['CODE_SIGN_STYLE']           = 'Automatic'
  config['DEVELOPMENT_TEAM']          = main_target
                                          .build_configuration_list
                                          .build_settings(config_name)
                                          .fetch('DEVELOPMENT_TEAM', '')
end

# ── Group for extension source files ─────────────────────────────────────────
groups_root = project.main_group
ext_group = groups_root.new_group(TARGET_NAME, TARGET_NAME)

# ── Source file ───────────────────────────────────────────────────────────────
swift_ref = ext_group.new_file('ShieldConfigurationExtension.swift')
ext_target.source_build_phase.add_file_reference(swift_ref)

# ── Info.plist (no build phase — just a file reference for Xcode to display) ─
ext_group.new_file('Info.plist')

# ── Entitlements ──────────────────────────────────────────────────────────────
ext_group.new_file("#{TARGET_NAME}.entitlements")

# ── Resource: app icon PNG ────────────────────────────────────────────────────
icon_ref = ext_group.new_file('AppIcon-1024.png')
resources_phase = ext_target.build_phases.find { |p| p.is_a?(Xcodeproj::Project::Object::PBXResourcesBuildPhase) }
unless resources_phase
  resources_phase = ext_target.new_build_phase(:resources)
end
resources_phase.add_file_reference(icon_ref)

# ── System frameworks ─────────────────────────────────────────────────────────
[
  'ShieldConfiguration',
  'ManagedSettings',
  'UIKit',
].each do |fw|
  fw_ref = project.frameworks_group.new_file(
    "System/Library/Frameworks/#{fw}.framework"
  )
  fw_ref.source_tree = 'SDKROOT'
  fw_ref.last_known_file_type = 'wrapper.framework'
  ext_target.frameworks_build_phase.add_file_reference(fw_ref)
end

# ── Embed extension in main app ──────────────────────────────────────────────
embed_phase = main_target.build_phases.find { |p|
  p.respond_to?(:name) && p.name == 'Embed Foundation Extensions'
}
unless embed_phase
  embed_phase = project.new(Xcodeproj::Project::Object::PBXCopyFilesBuildPhase)
  embed_phase.name = 'Embed Foundation Extensions'
  embed_phase.symbol_dst_subfolder_spec = :plug_ins
  main_target.build_phases << embed_phase
end

appex_ref = project.products_group.new_file("#{TARGET_NAME}.appex")
appex_ref.explicit_file_type = 'wrapper.app-extension'
appex_ref.source_tree = 'BUILT_PRODUCTS_DIR'
appex_ref.include_in_index = '0'

embed_file = embed_phase.add_file_reference(appex_ref)
embed_file.settings = { 'ATTRIBUTES' => ['RemoveHeadersOnCopy'] }

# ── Target dependency so main app builds extension first ─────────────────────
dep = project.new(Xcodeproj::Project::Object::PBXTargetDependency)
proxy = project.new(Xcodeproj::Project::Object::PBXContainerItemProxy)
proxy.container_portal = project.root_object.uuid
proxy.proxy_type       = '1'
proxy.remote_global_id_string = ext_target.uuid
proxy.remote_info      = TARGET_NAME
dep.target_proxy = proxy
dep.target       = ext_target
main_target.dependencies << dep

project.save
puts "✓ #{TARGET_NAME} target added and embedded in RemainFaithful"
