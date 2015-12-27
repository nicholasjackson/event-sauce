Before do |scenario|

end

After do |scenario|
  Registration.delete_all

  if scenario.failed?
    $FAILED = $FAILED + 1
  end

end
