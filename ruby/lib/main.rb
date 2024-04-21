require 'json'
require_relative 'evaluator'

if __FILE__ == $0
  dev_mode = false

  if dev_mode
    program = JSON.parse File.read '../test8.json'
  else
    data = ''
    while (input = gets)
      data += input
    end

    # Parse JSON data with potential errors due to strange file encodings
    program = nil
    loop do
      begin
        program = JSON.parse data
        break
      rescue JSON::ParserError => e
        data = data[1..-1]
      end
    end
  end

  puts (Evaluator.new program).evaluate_program
end