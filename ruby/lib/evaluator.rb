# frozen_string_literal: true

require_relative 'model'

class Evaluator
  def initialize(program)
    @program = program
  end

  def evaluate_program
    top_scope = Scope.new nil
    result = handle_body @program['body'], top_scope
    return result.str if result.str != ""
    new_final_value "function"
  end

  private

  def evaluate_expression(e)
    case e['type']
    when "Identifier"
      val = evaluate_identifier e
      return new_string_value new_error "unbound identifier" unless val
      return val
    when "Literal"
      new_string_value evaluate_literal e
    when "BinaryExpression"
      new_string_value evaluate_binary e
    when "UnaryExpression"
      new_string_value evaluate_unary e
    when "LogicalExpression"
      new_string_value evaluate_logical e
    when "ConditionalExpression"
      evaluate_conditional e
    when "FunctionExpression"
      new_function_value evaluate_function e
    when "CallExpression"
      evaluate_call e
    when "AssignmentExpression"
      new_string_value evaluate_assignment e
    else
      nil
    end
  end

  def evaluate_identifier(e)
    expr = get_identifier_value e['name'], e['scope']
    return nil unless expr
    expr
  end

  def evaluate_literal(e)
    new_final_value e['raw']
  end

  def evaluate_binary(e)
    e['left']['scope'] = e['scope']
    left = (evaluate_expression e['left']).str
    return left if has_error? left
    e['right']['scope'] = e['scope']
    right = (evaluate_expression e['right']).str
    return right if has_error? right
    left_numeric = get_number_from_value left
    return new_error "invalid binary type(s)" unless left_numeric
    right_numeric = get_number_from_value right
    return new_error "invalid binary type(s)" unless right_numeric
    if is_arithmetic_operator? e['operator']
      do_math left_numeric, right_numeric, e['operator']
    else
      do_comparison left_numeric, right_numeric, e['operator']
    end
  end

  def evaluate_unary(e)
    e['argument']['scope'] = e['scope']
    arg = (evaluate_expression e['argument']).str
    return arg if has_error? arg
    return new_error "invalid unary type" unless value_is_boolean? arg
    new_final_value bool_as_string get_bool_from_value !arg
  end

  def evaluate_logical(e)
    e['left']['scope'] = e['scope']
    left = (evaluate_expression e['left']).str
    return left if has_error? left
    e['right']['scope'] = e['scope']
    right = (evaluate_expression e['right']).str
    return right if has_error? right
    return new_error "invalid logical type(s)" unless (value_is_boolean? left) && (value_is_boolean? right)
    left_val = get_bool_from_value left
    right_val = get_bool_from_value right
    if e['operator'] == "||"
      result = left_val || right_val
    elsif e['operator'] == "&&"
      result = left_val && right_val
    end
    new_final_value bool_as_string result
  end

  def evaluate_conditional(e)
    e['test']['scope'] = e['scope']
    test = (evaluate_expression e['test']).str
    return new_string_value test if has_error? test
    return new_string_value new_error "invalid conditional type(s)" unless value_is_boolean? test
    if get_bool_from_value test
      e['consequent']['scope'] = e['scope']
      evaluate_expression e['consequent']
    else
      e['alternate']['scope'] = e['scope']
      evaluate_expression e['alternate']
    end
  end

  def evaluate_function(e)
    params = []
    e['params'].each do |param|
      params << (Parameter.new param['name'])
    end
    e['body']['body'].each do |line|
      unless line['argument'].nil?
        line['argument']['scope'] = Scope.new e['scope']
      end
    end
    Function.new params, e['body'], e['scope']
  end

  def evaluate_call(e)
    args = []
    e['arguments'].each do |arg|
      arg['scope'] = e['scope']
      result = evaluate_expression(arg)
      return result if has_error? result.str
      args << result
    end
    if e['callee']['name'].nil?
      e['callee']['scope'] = e['scope']
      result = evaluate_expression(e['callee'])
      return new_string_value new_error "not a function" unless result.str.nil?
      f = result.fun
    else
      val = get_identifier_value e['callee']['name'], e['scope']
      return new_string_value new_error "unbound identifier" unless val
      return new_string_value new_error "not a function" unless val.str.nil?
      f = val.fun
    end
    args.each_with_index do |arg, i|
      f.parameters[i].value = arg
    end
    execute_function f
  end

  def evaluate_assignment(e)
    unless get_identifier_value e['left']['name'], e['scope']
      return new_error "unbound identifier"
    end
    e['right']['scope'] = e['scope']
    value = evaluate_expression e['right']
    return value.str if has_error? value.str
    set_identifier_value e['left']['name'], value, e['scope']
    new_final_value "void"
  end

  def execute_function(function)
    param_scope = Scope.new function.scope
    function.parameters.each do |param|
      param_scope.variables[param.name] = param.value
    end
    bottom_scope = Scope.new param_scope
    handle_body function.body['body'], bottom_scope
  end

  def handle_body(body, scope)
    body.each_with_index do |statement, i|
      if statement['type'] == 'VariableDeclaration'
        err = handle_declarations statement['declarations'], scope
        return new_string_value err if err
      elsif statement['type'] == 'ReturnStatement'
        statement['argument']['scope'] = scope
        return evaluate_expression statement['argument']
      else
        statement['expression']['scope'] = scope
        value = evaluate_expression statement['expression']
        return value if has_error? value.str
        return value if !scope.parent && i == body.length - 1
      end
    end
  end

  def handle_declarations(declarations, scope)
    declarations.each do |declaration|
      name = declaration['id']['name']
      if declaration['init']['type'] == 'FunctionExpression'
        declaration['init']['scope'] = scope
        function = evaluate_function declaration['init']
        scope.variables[name] = new_function_value function
      else
        declaration['init']['scope'] = scope
        expr = evaluate_expression declaration['init']
        if has_error? expr.str
          return expr.str
        else
          scope.variables[name] = expr
        end
      end
    end
    nil
  end

  ### Simple helper methods ###

  def do_math(left, right, op)
    case op
    when '+'
      result = left + right
    when '-'
      result = left - right
    when '*'
      result = left * right
    when '/'
      if right == 0
        return new_error 'divide by zero'
      end
      result = left / right
    else
      return nil
    end
    new_final_value result.to_s
  end

  def do_comparison(left, right, op)
    case op
    when '=='
      result = left == right
    when '<'
      result = left < right
    when '>'
      result = right > left
    when '<='
      result = left <= right
    when '>='
      result = right >= left
    else
      return nil
    end
    new_final_value bool_as_string result
  end

  def is_arithmetic_operator?(op)
    op == '+' || op == '-' || op == '*' || op == '/'
  end

  def has_error?(str)
    str.include? 'error'
  end

  def value_is_boolean?(str)
    str.include? 'boolean'
  end

  def get_number_from_value(str)
    re = /\(value \(number (-?\d+)\)\)/

    match = str.match re
    if match.nil? || match.length < 2
      return nil
    end

    match[1].to_i # fixme: this will return 0 if not an integer
  end

  def get_bool_from_value(str)
    str.include? 'true'
  end

  def new_error(reason)
    "(error \"#{reason} banana\")"
  end

  def new_string_value(str)
    Value.new(str, nil)
  end

  def new_function_value(fun)
    Value.new(nil, fun)
  end

  def new_final_value(str)
    if str == 'true' || str == 'false'
      result = 'boolean ' + str
    elsif str == 'function' || str == 'void'
      result = str
    else
      if !str.to_i.is_a? Integer
        return new_error 'not a whole number'
      else
        result = 'number ' + str
      end
    end
    "(value (#{result}))"
  end

  def bool_as_string(bool)
    if bool
      'true'
    else
      'false'
    end
  end

  def get_identifier_value(id, scope)
    while scope
      value = scope.variables[id]
      if value
        return value
      else
        scope = scope.parent
      end
    end
  end

  def set_identifier_value(id, value, scope)
    while scope
      if scope.variables[id]
        scope.variables[id] = value
        return
      else
        scope = scope.parent
      end
    end
  end
end
