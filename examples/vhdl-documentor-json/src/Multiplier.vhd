library IEEE;

use IEEE.std_logic_1164.all;

entity Multiplier is

  port(
    iCLK : in  std_logic;
    i_A  : in  integer;
    i_B  : in  integer;
    o_P  : out integer
    );

end Multiplier;

architecture behavior of Multiplier is
begin

  process(iCLK, i_A, i_B)
  begin
    if rising_edge(iCLK) then
      o_P <= i_A * i_B;
    end if;
  end process;

end behavior;
