module Clipby

  class ReplacePlugin < Plugin

    def initialize
      @name = "Replace Plugin"

      @description <<-description
      Replace the letter e with ZZ from all data put into the system clipboard.
      description

      @author = "Dustin Willis Webber"
      @version = "0.1.0"
    end

    def self.run(type, data)
      log "Got Data!"
      return type, data.gsub!("e", "ZZ")
    end

  end

end
