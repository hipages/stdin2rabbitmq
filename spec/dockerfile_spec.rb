require 'dockerspec/serverspec'

describe "Container" do

  before :all do
    image = Docker::Image.get(ENV['DOCKER_IMAGE'])

    set :os, family: :alpine
    set :backend, :docker
    set :docker_image, image.id
  end

end
