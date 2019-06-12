# frozen_string_literal: true

require 'json'
require 'uri'
require 'rest-client'

# rubocop:disable Metrics/BlockLength
RSpec.describe 'Shorten URL API integrate test' do
  let(:base_url) { ENV['BASE_URL'] }
  let(:client) { RestClient }

  describe 'shorten API' do
    it 'response 200 with shortened url' do
      res = client.post(
        base_url,
        { origin: 'https://some.original.com' }.to_json,
        content_type: :json, accept: :json
      )
      expect(res.code).to eq 200

      json = JSON.parse(res.body, symbolize_names: true)
      expect(json.keys.to_set).to eq Set[:shorten]

      shorten_url = URI.parse(json[:shorten])
      expect(shorten_url.scheme).to eq 'https'
      expect(shorten_url.host).to eq URI.parse(base_url).host
      expect(shorten_url.path).to match %r{^\/\w{6}}
    end

    it 'response 400 when params invalid' do
      expect do
        client.post(
          base_url,
          {}.to_json,
          content_type: :json, accept: :json
        )
      end.to raise_error(RestClient::BadRequest)
    end
  end

  describe 'restore API' do
    it 'response 301 with valid token' do
      res = client.post(
        base_url,
        { origin: 'https://some.original.com' }.to_json,
        content_type: :json, accept: :json
      )
      expect(res.code).to eq 200

      json = JSON.parse(res.body, symbolize_names: true)
      expect(json.keys.to_set).to eq Set[:shorten]

      shorten_url = URI.parse(json[:shorten])

      # https://github.com/rest-client/rest-client#manually-following-redirection
      begin
        RestClient::Request.execute(
          method: :get,
          url: shorten_url.to_s,
          max_redirects: 0
        )
      rescue RestClient::MovedPermanently => e
        expect(e.response.code).to eq 301
        expect(e.response.headers[:location]).to eq 'https://some.original.com'
      end
    end

    it 'response 400 with invalid tokend' do
      expect do
        client.get(
          base_url + '/invalid_token',
        )
      end.to raise_error(RestClient::BadRequest)
    end

    it 'response 403 with empty token path' do
      expect { client.get(base_url) }.to raise_error(RestClient::Forbidden)
    end
  end
end
# rubocop:enable Metrics/BlockLength
