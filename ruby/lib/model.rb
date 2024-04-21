# frozen_string_literal: true

class Scope
  attr_accessor :parent
  attr_accessor :variables

  def initialize(parent)
    @parent = parent
    @variables = Hash.new
  end
end

class Function
  attr_accessor :parameters
  attr_accessor :body
  attr_accessor :scope

  def initialize(parameters, body, scope)
    @parameters = parameters
    @body = body
    @scope = scope
  end
end

class Value
  attr_accessor :str
  attr_accessor :fun

  def initialize(s, f)
    @str = s
    @fun = f
  end
end

class Parameter
  attr_accessor :name
  attr_accessor :value

  def initialize(name)
    @name = name
  end
end