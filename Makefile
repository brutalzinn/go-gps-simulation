# Define variables
ADB_URL_MACOSX := https://dl.google.com/android/repository/platform-tools-latest-darwin.zip?hl=pt-br

# Define install target for macOS
define INSTALL_ADB_MACOSX
  # Download ADB for macOS
  curl -L $(ADB_URL_MACOSX) > platform-tools-macos.zip
  # Extract the downloaded archive
  unzip platform-tools-macos.zip
  # Move ADB binaries to a specific location (modify as needed)
  mv platform-tools-macos/platform-tools/ adb/
  # Cleanup downloaded archive
  rm -rf platform-tools-macos.zip;
endef

# Installation target for macOS (call the macro)
install_adb: INSTALL_ADB_MACOSX

# Define a phony target to run the installation
.PHONY: all
all: install_adb
