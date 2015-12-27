Before do |scenario|

end

After do |scenario|
  Registration.delete_all

  if scenario.failed?
    $FAILED = $FAILED + 1
  end

end

begin
  e = $! # last exception
ensure
  puts 'arse' if $! != e
  raise e if $! != e
end
