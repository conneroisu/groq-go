library IEEE;

use IEEE.std_logic_1164.all;

entity mux2t1_N is

  generic(
    N : integer := 16
    );
  port(
    i_S  : in  std_logic;
    i_D0 : in  std_logic_vector(N-1 downto 0);
    i_D1 : in  std_logic_vector(N-1 downto 0);
    o_O  : out std_logic_vector(N-1 downto 0)
    );

end mux2t1_N;

architecture structural of mux2t1_N is

  component mux2t1 is
    port(
      i_S  : in  std_logic;
      i_D0 : in  std_logic;
      i_D1 : in  std_logic;
      o_O  : out std_logic
      );
  end component;

begin

  G_NBit_MUX : for i in 0 to N-1 generate
    MUXI : mux2t1 port map(
      i_S  => i_S,    -- All instances share the same select input.
      i_D0 => i_D0(i),  -- ith instance's data 0 input hooked up to ith data 0 input.
      i_D1 => i_D1(i),  -- ith instance's data 1 input hooked up to ith data 1 input.
      o_O  => o_O(i)  -- ith instance's data output hooked up to ith data output.
      );
  end generate G_NBit_MUX;

end structural;
