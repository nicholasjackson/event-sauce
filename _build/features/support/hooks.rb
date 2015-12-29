require 'pry'

Before do |scenario|
  Registration.delete_all
  clear_queue
end

After do |scenario|
  # Do something after each scenario.
  # The +scenario+ argument is optional, but
  # if you use it, you can inspect status with
  # the #failed?, #passed? and #exception methods.

  if scenario.failed?
    #binding.pry
  end

  Registration.delete_all
  clear_queue
end
