Before do |scenario|

end

After do |scenario|
  Registration.delete_all
end

at_exit do

end
