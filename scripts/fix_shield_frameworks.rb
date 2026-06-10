require 'xcodeproj'

PROJECT_PATH = File.expand_path('../../RemainFaithful.xcodeproj', __FILE__)
project = Xcodeproj::Project.open(PROJECT_PATH)

ext_target = project.targets.find { |t| t.name == 'RemainFaithfulShieldConfig' }
abort('RemainFaithfulShieldConfig target not found') unless ext_target

# Remove all existing framework references from the build phase
fw_phase = ext_target.frameworks_build_phase
fw_phase.files.dup.each { |f| fw_phase.remove_build_file(f) }

# Add the correct system frameworks
%w[ManagedSettings ManagedSettingsUI UIKit].each do |fw|
  ext_target.add_system_framework(fw)
end

project.save
puts "✓ Framework references updated for RemainFaithfulShieldConfig"
